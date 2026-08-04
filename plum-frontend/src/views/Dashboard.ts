import { defineComponent } from 'vue';
import { TimePeriod } from '@/dto/TimePeriod';
import { GroupingType } from '@/dto/GroupingType';
import type { RangeData } from '@/dto/Data';
import { ApiService } from '@/services/ApiService';
import HitsAndVisits from '@/components/HitsAndVisits.vue';
import Summary from '@/components/Summary.vue';
import Pages from '@/components/Pages.vue';
import Referers from '@/components/Referers.vue';
import BytesSent from '@/components/BytesSent.vue';
import BytesSentChart from '@/components/BytesSentChart.vue';
import StatusCodesChart from '@/components/StatusCodesChart.vue';
import StatusCodes from '@/components/StatusCodes.vue';

const apiService = new ApiService();

export default defineComponent({
    name: 'Dashboard',

    components: {
        HitsAndVisits,
        Summary,
        Pages,
        Referers,
        BytesSent,
        BytesSentChart,
        StatusCodesChart,
        StatusCodes,
    },

    setup() {
        // Exposed to the template so that it can reference the enum members.
        return {
            TimePeriod,
            GroupingType,
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

            data: [] as RangeData[],
            rangeData: [] as RangeData[],
            selectedRangeData: null as RangeData,

            timeoutID: null as ReturnType<typeof setTimeout>,
        };
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
            if (!this.updating) {
                this.reloadData();
            }
        },

        selectGroupingType(groupingType: GroupingType): void {
            this.selectedGroupingType = groupingType;
            if (!this.updating) {
                this.reloadData();
            }
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
            if (!this.updating) {
                this.reloadData();
            }
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
                    }
                });
        },

        reloadData(): void {
            this.updating = true;
            const timePeriod = this.selectedTimePeriod;
            const groupingType = this.selectedGroupingType;
            const selectedWebsite = this.selectedWebsite;
            apiService.getTimeRange(selectedWebsite, timePeriod, groupingType)
                .then(response => {
                    this.updating = false;
                    this.rangeData = response.data;
                    this.updateSelectedData();
                    this.scheduleReload(selectedWebsite, timePeriod, groupingType);
                });
        },

        scheduleReload(selectedWebsite: string, timePeriod: TimePeriod, groupingType: GroupingType): void {
            this.cancelReload();
            this.timeoutID = setTimeout(() => this.reloadLatestData(selectedWebsite, timePeriod, groupingType), 5000);
        },

        cancelReload(): void {
            if (this.timeoutID) {
                clearTimeout(this.timeoutID);
            }
        },

        reloadLatestData(selectedWebsite: string, timePeriod: TimePeriod, groupingType: GroupingType): void {
            this.updatingLatest = true;
            apiService.getTimePoint(selectedWebsite, timePeriod, groupingType)
                .then(response => {
                    this.updatingLatest = false;
                    if (selectedWebsite === this.selectedWebsite && timePeriod === this.selectedTimePeriod && groupingType === this.selectedGroupingType) {
                        this.updateLatestData(response.data);
                        this.scheduleReload(selectedWebsite, timePeriod, groupingType);
                    }
                })
                .catch(() => {
                    this.updatingLatest = false;
                    this.scheduleReload(selectedWebsite, timePeriod, groupingType);
                });
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
