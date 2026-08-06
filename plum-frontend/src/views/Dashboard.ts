import { defineComponent } from 'vue';
import { TimePeriod } from '@/dto/TimePeriod';
import { GroupingType } from '@/dto/GroupingType';
import { Category, CategoryLabels } from '@/dto/Category';
import type { RangeData } from '@/dto/Data';
import { FilterDimension, FilterLabels, filtersEqual, type Filter } from '@/dto/Filter';
import { ApiService } from '@/services/ApiService';
import CategoryTraffic from '@/components/CategoryTraffic.vue';
import HitsAndVisits from '@/components/HitsAndVisits.vue';
import Pages from '@/components/Pages.vue';
import Referers from '@/components/Referers.vue';
import BytesSent from '@/components/BytesSent.vue';
import BytesSentChart from '@/components/BytesSentChart.vue';
import StatusCodesChart from '@/components/StatusCodesChart.vue';
import StatusCodes from '@/components/StatusCodes.vue';

class ActiveFilter {
    dimension: FilterDimension;
    label: string;
    value: string;
}

const apiService = new ApiService();

export default defineComponent({
    name: 'Dashboard',

    components: {
        CategoryTraffic,
        HitsAndVisits,
        Pages,
        Referers,
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

            data: [] as RangeData[],
            rangeData: [] as RangeData[],
            selectedRangeData: null as RangeData,

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
                this.selectedRangeData ? this.selectedRangeData.time : '',
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
            this.reloadData();
        },

        selectGroupingType(groupingType: GroupingType): void {
            this.selectedGroupingType = groupingType;
            this.reloadData();
        },

        addFilter(dimension: FilterDimension, value: string): void {
            if (this.filter[dimension] === value) {
                return;
            }
            this.filter = {...this.filter, [dimension]: value};
            this.filterChanged();
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

        removeFilter(dimension: FilterDimension): void {
            const filter = {...this.filter};
            delete filter[dimension];
            this.filter = filter;
            this.filterChanged();
        },

        filterChanged(): void {
            this.selectedRangeData = null;
            this.reloadData();
        },

        selectData(index: number): void {
            if (index !== null && index !== undefined && index >= 0 && index < this.rangeData.length) {
                this.selectedRangeData = this.rangeData[index];
            } else {
                this.selectedRangeData = null;
            }
            this.updateSelectedData();
        },

        groupingTypeAvailable(groupingType: GroupingType): boolean {
            return this.getAvailableGroupingTypes(this.selectedTimePeriod)
                .some(v => v === groupingType);
        },

        selectWebsite(website: string): void {
            this.selectedWebsite = website;
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
            this.cancelReload();
            const timePeriod = this.selectedTimePeriod;
            const groupingType = this.selectedGroupingType;
            const selectedWebsite = this.selectedWebsite;
            const filter = this.filter;
            apiService.getTimeRange(selectedWebsite, timePeriod, groupingType, filter)
                .then(response => {
                    if (!this.parametersUnchanged(selectedWebsite, timePeriod, groupingType, filter)) {
                        return;
                    }
                    this.updating = false;
                    this.rangeData = response.data;
                    this.updateSelectedData();
                    this.scheduleReload(selectedWebsite, timePeriod, groupingType, filter);
                })
                .catch(() => {
                    if (this.parametersUnchanged(selectedWebsite, timePeriod, groupingType, filter)) {
                        this.updating = false;
                    }
                });
        },

        scheduleReload(selectedWebsite: string, timePeriod: TimePeriod, groupingType: GroupingType, filter: Filter): void {
            this.cancelReload();
            this.timeoutID = setTimeout(() => this.reloadLatestData(selectedWebsite, timePeriod, groupingType, filter), 5000);
        },

        cancelReload(): void {
            if (this.timeoutID) {
                clearTimeout(this.timeoutID);
            }
        },

        reloadLatestData(selectedWebsite: string, timePeriod: TimePeriod, groupingType: GroupingType, filter: Filter): void {
            this.updatingLatest = true;
            apiService.getTimePoint(selectedWebsite, timePeriod, groupingType, filter)
                .then(response => {
                    this.updatingLatest = false;
                    if (this.parametersUnchanged(selectedWebsite, timePeriod, groupingType, filter)) {
                        this.updateLatestData(response.data);
                        this.scheduleReload(selectedWebsite, timePeriod, groupingType, filter);
                    }
                })
                .catch(() => {
                    this.updatingLatest = false;
                    if (this.parametersUnchanged(selectedWebsite, timePeriod, groupingType, filter)) {
                        this.scheduleReload(selectedWebsite, timePeriod, groupingType, filter);
                    }
                });
        },

        parametersUnchanged(selectedWebsite: string, timePeriod: TimePeriod, groupingType: GroupingType, filter: Filter): boolean {
            return selectedWebsite === this.selectedWebsite
                && timePeriod === this.selectedTimePeriod
                && groupingType === this.selectedGroupingType
                && filtersEqual(filter, this.filter);
        },

        updateLatestData(rangeData: RangeData): void {
            if (this.rangeData.length > 0) {
                const data = Array.from(this.rangeData);
                const lastIndex = data.length - 1;
                if (rangeData.time === data[lastIndex].time) {
                    data[lastIndex] = rangeData;
                } else {
                    data.push(rangeData);
                    data.shift();
                }
                this.rangeData = data;
                this.updateSelectedData();
            }
        },

        updateSelectedData(): void {
            if (this.selectedRangeData) {
                this.data = [this.selectedRangeData];
            } else {
                this.data = this.rangeData;
            }
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
