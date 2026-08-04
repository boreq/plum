import type { Data, Dictionary, Metrics, RangeData } from '@/dto/Data';

export interface NamedMetrics extends Metrics {
    name: string;
}

export type MetricsSelector = (data: Data) => Dictionary<Metrics>;

export class MetricsService {

    total(data: RangeData[], selector: (data: Data) => number): number {
        if (!data) {
            return 0;
        }
        return data.reduce((acc, v) => acc + (v.data ? selector(v.data) : 0), 0);
    }

    group(data: RangeData[], selector: MetricsSelector): NamedMetrics[] {
        const rv: NamedMetrics[] = [];
        if (!data) {
            return rv;
        }
        for (const rangeData of data) {
            if (!rangeData.data) {
                continue;
            }
            const dimension = selector(rangeData.data);
            if (!dimension) {
                continue;
            }
            Object.entries(dimension).forEach(([name, metrics]) => {
                let row = rv.find(v => v.name === name);
                if (!row) {
                    row = {
                        name: name,
                        visits: 0,
                        hits: 0,
                        bytes: 0,
                    };
                    rv.push(row);
                }
                row.visits += metrics.visits;
                row.hits += metrics.hits;
                row.bytes += metrics.bytes;
            });
        }
        return rv;
    }

}
