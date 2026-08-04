import { defineComponent, type PropType } from 'vue';
import type { RangeData } from '@/dto/Data';
import { Align, type TableHeader, type TableRow } from '@/dto/Table';
import { UriService } from '@/services/UriService';
import { TextService } from '@/services/TextService';
import Table from '@/components/Table.vue';

class UriData {
    uri: string;
    visits: number;
    hits: number;
}

const uriService = new UriService();
const textService = new TextService();

export default defineComponent({
    name: 'Pages',

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
                        label: 'Page',
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
            const rows: UriData[] = [];
            for (const rangeData of this.data) {
                if (rangeData.data.uris) {
                    Object.entries(rangeData.data.uris).forEach(([uri, uriData]) => {
                        let row = rows.find(v => v.uri === uri);
                        if (!row) {
                            row = {
                                uri: uri,
                                visits: 0,
                                hits: 0,
                            };
                            rows.push(row);
                        }
                        row.visits += uriData.visits;
                        Object.entries(uriData.statuses).forEach(([, statusData]) => {
                            row.hits += statusData.hits;
                        });
                    });
                }
            }
            return this.toTableRows(rows);
        },
    },

    methods: {
        toTableRows(uriData: UriData[]): TableRow[] {
            const total: number = uriData.reduce((acc, v) => acc + v.visits, 0);
            return uriData
                .sort((a, b) => a.visits < b.visits ? 1 : -1)
                .filter(v => !uriService.isStaticResource(v.uri))
                .map(v => {
                    return {
                        data: [
                            v.uri,
                            textService.humanizeNumber(v.hits),
                            textService.humanizeNumber(v.visits),
                        ],
                        fraction: v.visits / total,
                    };
                });
        },
    },
});
