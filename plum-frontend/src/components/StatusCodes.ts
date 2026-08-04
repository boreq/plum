import { defineComponent, type PropType } from 'vue';
import type { Dictionary, RangeData } from '@/dto/Data';
import { Align, type TableHeader, type TableRow } from '@/dto/Table';
import { DataService } from '@/services/DataService';
import { TextService } from '@/services/TextService';
import Table from '@/components/Table.vue';
import TablePopup from '@/components/TablePopup.vue';

class StatusCodesData {
    status: string;
    uris: Dictionary<number>;
}

const dataService = new DataService();
const textService = new TextService();

export default defineComponent({
    name: 'StatusCodes',

    components: {
        Table,
        TablePopup,
    },

    props: {
        data: {
            type: Array as PropType<RangeData[]>,
            default: () => [],
        },
    },

    data() {
        return {
            modalRows: [] as TableRow[],
            modalTitle: '',
            displayModal: false,
        };
    },

    computed: {
        header(): TableHeader {
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
        },

        popupHeader(): TableHeader {
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
        },

        rows(): TableRow[] {
            if (!this.data) {
                return [];
            }
            return this.toTableRows(this.getStatusCodesData());
        },
    },

    methods: {
        clickRow(rowIndex: number): void {
            const row = this.getStatusCodesData()[rowIndex];
            this.modalRows = this.getChildren(row);
            this.modalTitle = textService.getHttpStatusText(row.status);
            this.displayModal = true;
        },

        getStatusCodesData(): StatusCodesData[] {
            const rows: StatusCodesData[] = [];
            for (const rangeData of this.data) {
                if (rangeData.data) {
                    Object.entries(dataService.getStatusMapping(rangeData.data))
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
        },

        toTableRows(statusCodesData: StatusCodesData[]): TableRow[] {
            const total: number = statusCodesData.reduce((acc, v) => acc + this.getTotal(v), 0);
            return statusCodesData
                .reduce<TableRow[]>((acc, v) => {
                    const statusTotal = this.getTotal(v);
                    acc.push({
                        data: [
                            textService.getHttpStatusText(v.status),
                            textService.humanizeNumber(statusTotal),
                        ],
                        fraction: statusTotal / total,
                    });
                    return acc;
                }, []);
        },

        getTotal(statusCodesData: StatusCodesData): number {
            return Object.entries(statusCodesData.uris)
                .reduce((acc, [, hits]) => {
                    return acc + hits;
                }, 0);
        },

        getChildren(v: StatusCodesData): TableRow[] {
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
        },
    },
});
