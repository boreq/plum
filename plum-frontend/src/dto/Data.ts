export interface Dictionary<T> {
    [key: string]: T;
}

export interface Metrics {
    visits: number;
    hits: number;
    bytes: number;
}

export interface UserAgentMetrics extends Metrics {
    browser: string;
}

export interface Data extends Metrics {
    categories: Dictionary<Metrics>;
    uris: Dictionary<Metrics>;
    statuses: Dictionary<Metrics>;
    referers: Dictionary<Metrics>;
    userAgents: Dictionary<UserAgentMetrics>;
}

export interface SeriesPoint extends Metrics {
    time: string;
    statuses: Dictionary<number>;
}

export interface RangeResult {
    summary: Data;
    series: SeriesPoint[];
}

export const emptyMetrics: Metrics = Object.freeze({
    visits: 0,
    hits: 0,
    bytes: 0,
});

export interface NamedMetrics extends Metrics {
    name: string;
}

export function namedMetrics(dimension: Dictionary<Metrics>): NamedMetrics[] {
    if (!dimension) {
        return [];
    }
    return Object.entries(dimension).map(([name, metrics]) => {
        return {
            name: name,
            visits: metrics.visits,
            hits: metrics.hits,
            bytes: metrics.bytes,
        };
    });
}
