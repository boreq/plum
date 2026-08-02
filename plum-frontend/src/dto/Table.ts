export class TableHeader {
    columns: TableHeaderColumn[];
}

export enum Align {
    Left = 'left',
    Right = 'right',
    Center = 'center',
}

export class TableHeaderColumn {
    label: string;
    width: string;
    align: Align;
}

export class TableRow {
    data: string[];
    fraction: number;
}
