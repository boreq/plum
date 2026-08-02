import { DateTime } from 'luxon';
import { TimePeriod } from '@/dto/TimePeriod';
import { GroupingType } from '@/dto/GroupingType';
import { RangeData } from '@/dto/Data';
import axios, { AxiosResponse } from 'axios'; // do not add { }, some webshit bs?

export class ApiService {

    getWebsites(): Promise<AxiosResponse<string[]>> {
        const url = `websites`;
        return axios.get<string[]>(process.env.VUE_APP_API_PREFIX + url);
    }

    getTimeRange(website: string, timePeriod: TimePeriod, groupingType: GroupingType): Promise<AxiosResponse<RangeData[]>> {
        const end = DateTime.utc();
        const start = end.minus(this.toLuxonRange(timePeriod));

        switch (groupingType) {
            case GroupingType.Hourly:
                return this.getHourly(website, start, end);
            case GroupingType.Daily:
                return this.getDaily(website, start, end);
            case GroupingType.Monthly:
                return this.getMonthly(website, start, end);
            default:
                throw new Error('not implemented');
        }
    }

    getTimePoint(website: string, timePeriod: TimePeriod, groupingType: GroupingType): Promise<AxiosResponse<RangeData>> {
        const now = DateTime.utc();

        switch (groupingType) {
            case GroupingType.Hourly:
                return this.getHour(website, now);
            case GroupingType.Daily:
                return this.getDay(website, now);
            case GroupingType.Monthly:
                return this.getMonth(website, now);
            default:
                throw new Error('not implemented');
        }
    }

    private getHour(website: string, t: DateTime): Promise<AxiosResponse<RangeData>> {
        const url = `hour/${t.year}/${t.month}/${t.day}/${t.hour}.json`;
        return axios.get<RangeData>(process.env.VUE_APP_API_PREFIX + this.websiteUrl(website) + url);
    }

    private getDay(website: string, t: DateTime): Promise<AxiosResponse<RangeData>> {
        const url = `day/${t.year}/${t.month}/${t.day}.json`;
        return axios.get<RangeData>(process.env.VUE_APP_API_PREFIX + this.websiteUrl(website) + url);
    }

    private getMonth(website: string, t: DateTime): Promise<AxiosResponse<RangeData>> {
        const url = `month/${t.year}/${t.month}.json`;
        return axios.get<RangeData>(process.env.VUE_APP_API_PREFIX + this.websiteUrl(website) + url);
    }

    private getHourly(website: string, from: DateTime, to: DateTime): Promise<AxiosResponse<RangeData[]>> {
        const url = `range/hourly/` +
            `${from.year}/${from.month}/${from.day}/${from.hour}/` +
            `${to.year}/${to.month}/${to.day}/${to.hour}.json`;
        return axios.get<RangeData[]>(process.env.VUE_APP_API_PREFIX + this.websiteUrl(website) + url);
    }

    private getDaily(website: string, from: DateTime, to: DateTime): Promise<AxiosResponse<RangeData[]>> {
        const url = `range/daily/` +
            `${from.year}/${from.month}/${from.day}/` +
            `${to.year}/${to.month}/${to.day}.json`;
        return axios.get<RangeData[]>(process.env.VUE_APP_API_PREFIX + this.websiteUrl(website) + url);
    }

    private getMonthly(website: string, from: DateTime, to: DateTime): Promise<AxiosResponse<RangeData[]>> {
        const url = `range/monthly/` +
            `${from.year}/${from.month}/` +
            `${to.year}/${to.month}.json`;
        return axios.get<RangeData[]>(process.env.VUE_APP_API_PREFIX + this.websiteUrl(website) + url);
    }

    private websiteUrl(website: string): string {
        return 'websites/' + encodeURIComponent(website) + '/';
    }

    private toLuxonRange(timePeriod: TimePeriod): any {
        switch (timePeriod) {
            case TimePeriod.Day:
                return {days: 1};
            case TimePeriod.Week:
                return {days: 6};
            case TimePeriod.Month:
                return {months: 1};
            case TimePeriod.Year:
                return {years: 1};
            default:
                return null;
        }
    }

}
