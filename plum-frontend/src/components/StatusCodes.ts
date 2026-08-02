import { Component, Prop, Vue } from 'vue-property-decorator';
import { Dictionary, RangeData } from '@/dto/Data';
import { Align, TableHeader, TableRow } from '@/dto/Table';
import { DataService } from '@/services/DataService';
import { TextService } from '@/services/TextService';
import Table from '@/components/Table.vue';
import TablePopup from '@/components/TablePopup.vue';

class StatusCodesData {
    status: string;
    uris: Dictionary<number>;
}

@Component({
    components: {
        Table,
        TablePopup,
    },
})
export default class StatusCodes extends Vue {

    @Prop()
    data: RangeData[];

    modalRows: TableRow[] = [];
    modalTitle: string;
    displayModal = false;

    private dataService = new DataService();
    private textService = new TextService();

    get header(): TableHeader {
        return {
            columns: [
                {
                    label: 'Status',
                    width: null,
                    align: Align.Left,
                },
                {
                    label: 'Hits',
                    width: '100px',
                    align: Align.Right,
                },
            ],
        };
    }

    get popupHeader(): TableHeader {
        return {
            columns: [
                {
                    label: 'Resource',
                    width: null,
                    align: Align.Left,
                },
                {
                    label: 'Hits',
                    width: '100px',
                    align: Align.Right,
                },
            ],
        };
    }

    get rows(): TableRow[] {
        if (!this.data) {
            return [];
        }
        const rows = this.getStatusCodesData();
        return this.toTableRows(rows);
    }

    clickRow(rowIndex: number): void {
        const row = this.getStatusCodesData()[rowIndex];
        const children = this.getChildren(row);
        this.modalRows = children;
        this.modalTitle = this.textService.getHttpStatusText(row.status);
        this.displayModal = true;
    }

    private getStatusCodesData(): StatusCodesData[] {
        const rows: StatusCodesData[] = [];
        for (const rangeData of this.data) {
            if (rangeData.data) {
                Object.entries(this.dataService.getStatusMapping(rangeData.data))
                    .forEach(([status, uriMap]) => {
                        let row = rows.find(v => v.status === status);
                        if (!row) {
                            row = {
                                status: status,
                                uris: {},
                            };
                            rows.push(row);
                        }
                        Object.entries(uriMap).forEach(([uri, hits]) => {
                            row.uris[uri] = (row.uris[uri] || 0) + hits;
                        });
                    });
            }
        }
        return rows.sort((a, b) => a.status > b.status ? 1 : -1);
    }

    private toTableRows(statusCodesData: StatusCodesData[]): TableRow[] {
        const total: number = statusCodesData.reduce((acc, v) => acc + this.getTotal(v), 0);
        return statusCodesData
            .reduce((acc, v) => {
                const statusTotal = this.getTotal(v);
                acc.push({
                    data: [
                        this.textService.getHttpStatusText(v.status),
                        this.textService.humanizeNumber(statusTotal),
                    ],
                    fraction: statusTotal / total,
                });
                return acc;
            }, []);
    }

    private getTotal(statusCodesData: StatusCodesData): number {
        return Object.entries(statusCodesData.uris)
            .reduce((acc, [uri, hits]) => {
                return acc + hits;
            }, 0);
    }

    private getChildren(v: StatusCodesData): TableRow[] {
        const total: number = this.getTotal(v);

        return Object.entries(v.uris)
            .map(([uri, hits]) => {
                return {
                    data: [
                        uri,
                        hits.toString(),
                    ],
                    fraction: hits / total,
                };
            })
            .sort((a, b) => Number(a.data[1]) < Number(b.data[1]) ? 1 : -1);
    }

}
