import { Component, Vue } from 'vue-property-decorator';
import { TimePeriod } from '@/dto/TimePeriod';
import { GroupingType } from '@/dto/GroupingType';
import { RangeData } from '@/dto/Data';
import { ApiService } from '@/services/ApiService';
import HitsAndVisits from '@/components/HitsAndVisits.vue';
import Summary from '@/components/Summary.vue';
import Pages from '@/components/Pages.vue';
import Referers from '@/components/Referers.vue';
import BytesSent from '@/components/BytesSent.vue';
import BytesSentChart from '@/components/BytesSentChart.vue';
import StatusCodesChart from '@/components/StatusCodesChart.vue';
import StatusCodes from '@/components/StatusCodes.vue';


@Component({
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
})
export default class Dashboard extends Vue {

    TimePeriod = TimePeriod;
    GroupingType = GroupingType;

    selectedTimePeriod: TimePeriod = TimePeriod.Day;
    selectedGroupingType: GroupingType = GroupingType.Hourly;

    updating: boolean = false;
    updatingLatest: boolean = false;

    websites: string[] = [];
    selectedWebsite: string = '';

    data: RangeData[] = [];
    rangeData: RangeData[] = [];
    selectedRangeData: RangeData = null;

    private apiService = new ApiService();
    private timeoutID: number;


    created(): void {
        this.loadWebsites();
    }

    beforeDestroy(): void {
        this.cancelReload();
    }

    selectTimePeriod(timePeriod: TimePeriod): void {
        this.selectedTimePeriod = timePeriod;
        this.selectAppropriateGroupingType(timePeriod);
        if (!this.updating) {
            this.reloadData();
        }
    }

    selectGroupingType(groupingType: GroupingType): void {
        this.selectedGroupingType = groupingType;
        if (!this.updating) {
            this.reloadData();
        }
    }

    selectData(index: number): void {
        if (index !== null && index !== undefined && index >= 0 && index < this.rangeData.length) {
            this.selectedRangeData = this.rangeData[index];
        } else {
            this.selectedRangeData = null;
        }
        this.updateSelectedData();
    }

    groupingTypeAvailable(groupingType: GroupingType): boolean {
        return this.getAvailableGroupingTypes(this.selectedTimePeriod)
            .some(v => v === groupingType);
    }

    selectWebsite(website: string): void {
        this.selectedWebsite = website;
        if (!this.updating) {
            this.reloadData();
        }
    }

    private selectAppropriateGroupingType(timePeriod: TimePeriod): void {
        switch (this.selectedTimePeriod) {
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
    }

    private loadWebsites(): void {
        this.updating = true;
        this.apiService.getWebsites()
            .then(response => {
                this.websites = response.data;
                if (this.websites && this.websites.length > 0) {
                    this.selectedWebsite = this.websites[0];
                    this.reloadData();
                }
            });
    }

    private reloadData(): void {
        this.updating = true;
        const timePeriod = this.selectedTimePeriod;
        const groupingType = this.selectedGroupingType;
        const selectedWebsite = this.selectedWebsite;
        this.apiService.getTimeRange(selectedWebsite, timePeriod, groupingType)
            .then(response => {
                this.updating = false;
                this.rangeData = response.data;
                this.updateSelectedData();
                this.scheduleReload(selectedWebsite, timePeriod, groupingType);
            });
    }

    private scheduleReload(selectedWebsite: string, timePeriod: TimePeriod, groupingType: GroupingType): void {
        this.cancelReload();
        this.timeoutID = setTimeout(() => this.reloadLatestData(selectedWebsite, timePeriod, groupingType), 5000);
    }

    private cancelReload(): void {
        if (this.timeoutID) {
            clearTimeout(this.timeoutID);
        }
    }

    private reloadLatestData(selectedWebsite: string, timePeriod: TimePeriod, groupingType: GroupingType): void {
        this.updatingLatest = true;
        this.apiService.getTimePoint(selectedWebsite, timePeriod, groupingType)
            .then(response => {
                this.updatingLatest = false;
                if (selectedWebsite === this.selectedWebsite && timePeriod === this.selectedTimePeriod && groupingType === this.selectedGroupingType) {
                    this.updateLatestData(response.data);
                    this.scheduleReload(selectedWebsite, timePeriod, groupingType);
                }
            })
            .catch(error => {
                this.updatingLatest = false;
                this.scheduleReload(selectedWebsite, timePeriod, groupingType);
            });
    }

    private updateLatestData(rangeData: RangeData): void {
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
    }

    private updateSelectedData(): void {
        if (this.selectedRangeData) {
            this.data = [this.selectedRangeData];
        } else {
            this.data = this.rangeData;
        }
    }

    private getAvailableGroupingTypes(timePeriod: TimePeriod): GroupingType[] {
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
    }

}
