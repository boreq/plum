import { defineComponent, type PropType } from 'vue';
import type { RangeData } from '@/dto/Data';
import { Align, type TableHeader, type TableRow } from '@/dto/Table';
import { TextService } from '@/services/TextService';
import Table from '@/components/Table.vue';

class RefererData {
    referer: string;
    visits: number;
    hits: number;
}

const textService = new TextService();

export default defineComponent({
    name: 'Referers',

    components: {
        Table,
    },

    props: {
        data: {
            type: Array as PropType<RangeData[]>,
            default: () => [],
        },
    },

    computed: {
        header(): TableHeader {
            return {
                columns: [
                    {
                        label: 'Referer',
                        width: null,
                        align: Align.Left,
                    },
                    {
                        label: 'Hits',
                        width: '60px',
                        align: Align.Right,
                    },
                    {
                        label: 'Visits',
                        width: '60px',
                        align: Align.Right,
                    },
                ],
            };
        },

        rows(): TableRow[] {
            if (!this.data) {
                return [];
            }
            const rows: RefererData[] = [];
            for (const rangeData of this.data) {
                if (rangeData.data.referers) {
                    Object.entries(rangeData.data.referers).forEach(([referer, refererData]) => {
                        let row = rows.find(v => v.referer === referer);
                        if (!row) {
                            row = {
                                referer: referer,
                                visits: 0,
                                hits: 0,
                            };
                            rows.push(row);
                        }
                        row.visits += refererData.visits;
                        row.hits += refererData.hits;
                    });
                }
            }
            return this.toTableRows(rows);
        },
    },

    methods: {
        toTableRows(refererData: RefererData[]): TableRow[] {
            const total: number = refererData.reduce((acc, v) => acc + v.visits, 0);
            return refererData
                .sort((a, b) => a.visits < b.visits ? 1 : -1)
                .map(v => {
                    return {
                        data: [
                            v.referer,
                            textService.humanizeNumber(v.hits),
                            textService.humanizeNumber(v.visits),
                        ],
                        fraction: v.visits / total,
                    };
                });
        },
    },
});
