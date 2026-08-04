export interface Dictionary<T> {
    [key: string]: T;
}

export interface Metrics {
    visits: number;
    hits: number;
    bytes: number;
}

export interface Data extends Metrics {
    uris: Dictionary<Metrics>;
    statuses: Dictionary<Metrics>;
    referers: Dictionary<Metrics>;
}

export interface RangeData {
    time: string;
    data: Data;
}
