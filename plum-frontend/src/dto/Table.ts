export class TableHeader {
    columns: TableHeaderColumn[];
}

export enum Align {
    Left = 'left',
    Right = 'right',
    Center = 'center',
}

export enum SortDirection {
    Ascending = 'ascending',
    Descending = 'descending',
}

export type TableValue = string | number;

export class TableHeaderColumn {
    label: string;
    width: string;
    align: Align;
    sortable?: boolean;
    format?: (value: TableValue) => string;
}

export class TableRow {
    data: TableValue[];
    fraction: number;
    icon?: string;
    iconTitle?: string;
}
