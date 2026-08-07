import { DateTime } from 'luxon';
import type { DurationLike } from 'luxon';
import { TimePeriod } from '@/dto/TimePeriod';
import { GroupingType } from '@/dto/GroupingType';
import type { Data, RangeResult } from '@/dto/Data';
import type { Filter } from '@/dto/Filter';
import axios from 'axios';
import type { AxiosResponse } from 'axios';

const API_PREFIX = '/api/';

export class ApiService {

    getWebsites(): Promise<AxiosResponse<string[]>> {
        const url = `websites`;
        return axios.get<string[]>(API_PREFIX + url);
    }

    getTimeRange(website: string, timePeriod: TimePeriod, groupingType: GroupingType, filter: Filter): Promise<AxiosResponse<RangeResult>> {
        const end = DateTime.utc();
        const start = end.minus(this.toLuxonRange(timePeriod));

        switch (groupingType) {
            case GroupingType.Hourly:
                return this.getHourly(website, start, end, filter);
            case GroupingType.Daily:
                return this.getDaily(website, start, end, filter);
            case GroupingType.Monthly:
                return this.getMonthly(website, start, end, filter);
            default:
                throw new Error('not implemented');
        }
    }

    getTimePoint(website: string, groupingType: GroupingType, time: string, filter: Filter): Promise<AxiosResponse<Data>> {
        const t = DateTime.fromISO(time).toUTC();

        switch (groupingType) {
            case GroupingType.Hourly:
                return this.getHour(website, t, filter);
            case GroupingType.Daily:
                return this.getDay(website, t, filter);
            case GroupingType.Monthly:
                return this.getMonth(website, t, filter);
            default:
                throw new Error('not implemented');
        }
    }

    private getHour(website: string, t: DateTime, filter: Filter): Promise<AxiosResponse<Data>> {
        const url = `hour/${t.year}/${t.month}/${t.day}/${t.hour}.json`;
        return axios.get<Data>(API_PREFIX + this.websiteUrl(website) + url, {params: filter});
    }

    private getDay(website: string, t: DateTime, filter: Filter): Promise<AxiosResponse<Data>> {
        const url = `day/${t.year}/${t.month}/${t.day}.json`;
        return axios.get<Data>(API_PREFIX + this.websiteUrl(website) + url, {params: filter});
    }

    private getMonth(website: string, t: DateTime, filter: Filter): Promise<AxiosResponse<Data>> {
        const url = `month/${t.year}/${t.month}.json`;
        return axios.get<Data>(API_PREFIX + this.websiteUrl(website) + url, {params: filter});
    }

    private getHourly(website: string, from: DateTime, to: DateTime, filter: Filter): Promise<AxiosResponse<RangeResult>> {
        const url = `range/hourly/` +
            `${from.year}/${from.month}/${from.day}/${from.hour}/` +
            `${to.year}/${to.month}/${to.day}/${to.hour}.json`;
        return axios.get<RangeResult>(API_PREFIX + this.websiteUrl(website) + url, {params: filter});
    }

    private getDaily(website: string, from: DateTime, to: DateTime, filter: Filter): Promise<AxiosResponse<RangeResult>> {
        const url = `range/daily/` +
            `${from.year}/${from.month}/${from.day}/` +
            `${to.year}/${to.month}/${to.day}.json`;
        return axios.get<RangeResult>(API_PREFIX + this.websiteUrl(website) + url, {params: filter});
    }

    private getMonthly(website: string, from: DateTime, to: DateTime, filter: Filter): Promise<AxiosResponse<RangeResult>> {
        const url = `range/monthly/` +
            `${from.year}/${from.month}/` +
            `${to.year}/${to.month}.json`;
        return axios.get<RangeResult>(API_PREFIX + this.websiteUrl(website) + url, {params: filter});
    }

    private websiteUrl(website: string): string {
        return 'websites/' + encodeURIComponent(website) + '/';
    }

    private toLuxonRange(timePeriod: TimePeriod): DurationLike {
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
