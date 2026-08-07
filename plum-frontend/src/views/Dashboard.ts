import { defineComponent } from 'vue';
import { TimePeriod } from '@/dto/TimePeriod';
import { GroupingType } from '@/dto/GroupingType';
import { Category, CategoryLabels } from '@/dto/Category';
import type { Data, SeriesPoint } from '@/dto/Data';
import { FilterDimension, FilterLabels, filtersEqual, type Filter } from '@/dto/Filter';
import { ApiService } from '@/services/ApiService';
import { TextService } from '@/services/TextService';
import CategoryTraffic from '@/components/CategoryTraffic.vue';
import HitsAndVisits from '@/components/HitsAndVisits.vue';
import Pages from '@/components/Pages.vue';
import Referers from '@/components/Referers.vue';
import UserAgents from '@/components/UserAgents.vue';
import BytesSent from '@/components/BytesSent.vue';
import BytesSentChart from '@/components/BytesSentChart.vue';
import StatusCodesChart from '@/components/StatusCodesChart.vue';
import StatusCodes from '@/components/StatusCodes.vue';

class ActiveFilter {
    dimension: FilterDimension;
    label: string;
    value: string;
}

class Parameters {
    website: string;
    timePeriod: TimePeriod;
    groupingType: GroupingType;
    filter: Filter;
    selectedTime: string;
}

const apiService = new ApiService();
const textService = new TextService();

const RELOAD_INTERVAL = 5000;

export default defineComponent({
    name: 'Dashboard',

    components: {
        CategoryTraffic,
        HitsAndVisits,
        Pages,
        Referers,
        UserAgents,
        BytesSent,
        BytesSentChart,
        StatusCodesChart,
        StatusCodes,
    },

    setup() {
        return {
            TimePeriod,
            GroupingType,
            Category,
        };
    },

    data() {
        return {
            selectedTimePeriod: TimePeriod.Day,
            selectedGroupingType: GroupingType.Hourly,

            updating: false,
            updatingLatest: false,

            websites: [] as string[],
            selectedWebsite: '',

            filter: {} as Filter,
            selectedTime: null as string,

            rangeSummary: null as Data,
            pointSummary: null as Data,
            series: [] as SeriesPoint[],

            timeoutID: null as ReturnType<typeof setTimeout>,
        };
    },

    computed: {
        activeFilters(): ActiveFilter[] {
            return Object.values(FilterDimension)
                .filter(dimension => this.filter[dimension])
                .map(dimension => {
                    return {
                        dimension: dimension,
                        label: FilterLabels[dimension],
                        value: this.filterValueLabel(dimension, this.filter[dimension]),
                    };
                });
        },

        selectedCategory(): Category {
            return this.filter[FilterDimension.Category] as Category;
        },

        summary(): Data {
            return this.selectedTime ? this.pointSummary : this.rangeSummary;
        },

        selectedTimeLabel(): string {
            return this.selectedTime ? textService.formatDate(this.selectedTime, this.selectedGroupingType) : '';
        },

        selectedIndex(): number {
            if (!this.selectedTime) {
                return null;
            }
            const index = this.series.findIndex(v => v.time === this.selectedTime);
            return index >= 0 ? index : null;
        },

        // Changes whenever a different set of data is loaded. Used as a key so
        // that the tables are recreated and their pagination is reset instead
        // of displaying a page which makes no sense for the new data. Periodic
        // reloads of the latest data point don't affect this value so that they
        // don't interrupt the user.
        tableKey(): string {
            const filter = Object.values(FilterDimension)
                .map(dimension => `${dimension}=${this.filter[dimension] || ''}`)
                .join('&');
            return [
                this.selectedWebsite,
                this.selectedTimePeriod,
                this.selectedGroupingType,
                this.selectedTime || '',
                filter,
            ].join('|');
        },
    },

    created(): void {
        this.loadWebsites();
    },

    beforeUnmount(): void {
        this.cancelReload();
    },

    methods: {
        selectTimePeriod(timePeriod: TimePeriod): void {
            this.selectedTimePeriod = timePeriod;
            this.selectAppropriateGroupingType(timePeriod);
            this.selectedTime = null;
            this.reloadData();
        },

        selectGroupingType(groupingType: GroupingType): void {
            this.selectedGroupingType = groupingType;
            this.selectedTime = null;
            this.reloadData();
        },

        addFilter(dimension: FilterDimension, value: string): void {
            if (this.filter[dimension] === value) {
                return;
            }
            this.filter = {...this.filter, [dimension]: value};
            this.filterChanged();
        },

        removeFilter(dimension: FilterDimension): void {
            const filter = {...this.filter};
            delete filter[dimension];
            this.filter = filter;
            this.filterChanged();
        },

        filterChanged(): void {
            this.selectedTime = null;
            this.reloadData();
        },

        categoryChecked(category: Category): boolean {
            return !this.selectedCategory || this.selectedCategory === category;
        },

        toggleCategory(category: Category): void {
            if (this.selectedCategory === category) {
                this.removeFilter(FilterDimension.Category);
            } else {
                this.addFilter(FilterDimension.Category, category);
            }
        },

        filterValueLabel(dimension: FilterDimension, value: string): string {
            if (dimension === FilterDimension.Category) {
                return CategoryLabels[value as Category] || value;
            }
            return value;
        },

        // Only the summary depends on the selected data point so the series
        // doesn't have to be loaded again.
        selectData(index: number): void {
            const clicked = index !== null && index !== undefined && index >= 0 && index < this.series.length
                ? this.series[index].time
                : null;

            const time = clicked === this.selectedTime ? null : clicked;
            if (time === this.selectedTime) {
                return;
            }

            this.selectedTime = time;
            this.pointSummary = null;
            this.updating = true;
            this.cancelReload();

            const parameters = this.currentParameters();
            this.loadPoint(parameters).then(() => this.loadFinished(parameters));
        },

        groupingTypeAvailable(groupingType: GroupingType): boolean {
            return this.getAvailableGroupingTypes(this.selectedTimePeriod)
                .some(v => v === groupingType);
        },

        selectWebsite(website: string): void {
            this.selectedWebsite = website;
            this.filter = {};
            this.selectedTime = null;
            this.reloadData();
        },

        selectAppropriateGroupingType(timePeriod: TimePeriod): void {
            switch (timePeriod) {
                case TimePeriod.Day:
                    this.selectedGroupingType = GroupingType.Hourly;
                    break;
                case TimePeriod.Week:
                    this.selectedGroupingType = GroupingType.Daily;
                    break;
                case TimePeriod.Month:
                    this.selectedGroupingType = GroupingType.Daily;
                    break;
                case TimePeriod.Year:
                    this.selectedGroupingType = GroupingType.Monthly;
                    break;
                default:
                    throw new Error('not implemented');
            }
        },

        loadWebsites(): void {
            this.updating = true;
            apiService.getWebsites()
                .then(response => {
                    this.websites = response.data;
                    if (this.websites && this.websites.length > 0) {
                        this.selectedWebsite = this.websites[0];
                        this.reloadData();
                    } else {
                        this.updating = false;
                    }
                })
                .catch(() => {
                    this.updating = false;
                });
        },

        reloadData(): void {
            this.updating = true;
            this.loadData();
        },

        reloadLatestData(): void {
            this.updatingLatest = true;
            this.loadData();
        },

        loadData(): void {
            this.cancelReload();

            const parameters = this.currentParameters();
            Promise.all([this.loadRange(parameters), this.loadPoint(parameters)])
                .then(() => this.loadFinished(parameters));
        },

        loadRange(parameters: Parameters): Promise<void> {
            return apiService.getTimeRange(parameters.website, parameters.timePeriod, parameters.groupingType, parameters.filter)
                .then(response => {
                    if (this.parametersUnchanged(parameters)) {
                        this.rangeSummary = response.data.summary;
                        this.series = response.data.series || [];
                    }
                })
                .catch(() => undefined);
        },

        loadPoint(parameters: Parameters): Promise<void> {
            if (!parameters.selectedTime) {
                this.pointSummary = null;
                return Promise.resolve();
            }

            return apiService.getTimePoint(parameters.website, parameters.groupingType, parameters.selectedTime, parameters.filter)
                .then(response => {
                    if (this.parametersUnchanged(parameters)) {
                        this.pointSummary = response.data.data;
                    }
                })
                .catch(() => undefined);
        },

        loadFinished(parameters: Parameters): void {
            if (!this.parametersUnchanged(parameters)) {
                return;
            }
            this.updating = false;
            this.updatingLatest = false;
            this.scheduleReload(parameters);
        },

        scheduleReload(parameters: Parameters): void {
            this.cancelReload();
            this.timeoutID = setTimeout(() => {
                if (this.parametersUnchanged(parameters)) {
                    this.reloadLatestData();
                }
            }, RELOAD_INTERVAL);
        },

        cancelReload(): void {
            if (this.timeoutID) {
                clearTimeout(this.timeoutID);
                this.timeoutID = null;
            }
        },

        currentParameters(): Parameters {
            return {
                website: this.selectedWebsite,
                timePeriod: this.selectedTimePeriod,
                groupingType: this.selectedGroupingType,
                filter: this.filter,
                selectedTime: this.selectedTime,
            };
        },

        parametersUnchanged(parameters: Parameters): boolean {
            return parameters.website === this.selectedWebsite
                && parameters.timePeriod === this.selectedTimePeriod
                && parameters.groupingType === this.selectedGroupingType
                && parameters.selectedTime === this.selectedTime
                && filtersEqual(parameters.filter, this.filter);
        },

        getAvailableGroupingTypes(timePeriod: TimePeriod): GroupingType[] {
            switch (timePeriod) {
                case TimePeriod.Day:
                    return [GroupingType.Hourly];
                case TimePeriod.Week:
                    return [GroupingType.Hourly, GroupingType.Daily];
                case TimePeriod.Month:
                    return [GroupingType.Hourly, GroupingType.Daily];
                case TimePeriod.Year:
                    return [GroupingType.Daily, GroupingType.Monthly];
                default:
                    throw new Error('not implemented');
            }
        },
    },
});
