import { defineComponent, type PropType } from 'vue';
import { namedMetrics, type Data, type NamedMetrics } from '@/dto/Data';
import { Align, type TableHeader, type TableRow } from '@/dto/Table';
import { FilterDimension } from '@/dto/Filter';
import { TextService } from '@/services/TextService';
import Table from '@/components/Table.vue';

const textService = new TextService();

export default defineComponent({
    name: 'Referers',

    components: {
        Table,
    },

    props: {
        data: {
            type: Object as PropType<Data>,
            default: null,
        },
    },

    emits: ['filter'],

    computed: {
        header(): TableHeader {
            return {
                columns: [
                    {
                        label: 'Referer',
                        width: null,
                        align: Align.Left,
                        sortable: false,
                    },
                    {
                        label: 'Hits',
                        width: '60px',
                        align: Align.Right,
                        format: v => textService.humanizeNumber(v as number),
                    },
                    {
                        label: 'Visits',
                        width: '60px',
                        align: Align.Right,
                        format: v => textService.humanizeNumber(v as number),
                    },
                ],
            };
        },

        referers(): NamedMetrics[] {
            return namedMetrics(this.data?.referers)
                .sort((a, b) => a.visits < b.visits ? 1 : -1);
        },

        rows(): TableRow[] {
            const total: number = this.referers.reduce((acc, v) => acc + v.visits, 0);
            return this.referers.map(v => {
                return {
                    data: [
                        v.name,
                        v.hits,
                        v.visits,
                    ],
                    fraction: total ? v.visits / total : 0,
                };
            });
        },
    },

    methods: {
        clickRow(rowIndex: number): void {
            this.$emit('filter', FilterDimension.Referer, this.referers[rowIndex].name);
        },
    },
});
