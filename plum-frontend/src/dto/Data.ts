export interface Dictionary<T> {
    [key: string]: T;
}

export class RangeData {
    time: string;
    data: Data;
}

export class Data {
    referers: Dictionary<RefererData>;
    uris: Dictionary<UriData>;
    visits: number;
}

export class RefererData {
    visits: number;
    hits: number;
}

export class UriData {
    visits: number;
    bytes: number;
    statuses: Dictionary<StatusData>;
}

export class StatusData {
    hits: number;
}
