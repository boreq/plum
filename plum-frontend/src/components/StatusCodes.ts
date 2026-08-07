import { defineComponent, type PropType } from 'vue';
import { namedMetrics, type Data, type NamedMetrics } from '@/dto/Data';
import { Align, type TableHeader, type TableRow } from '@/dto/Table';
import { FilterDimension } from '@/dto/Filter';
import { TextService } from '@/services/TextService';
import Table from '@/components/Table.vue';

const textService = new TextService();

export default defineComponent({
    name: 'StatusCodes',

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

        statusCodes(): NamedMetrics[] {
            return namedMetrics(this.data?.statuses)
                .sort((a, b) => a.name > b.name ? 1 : -1);
        },

        rows(): TableRow[] {
            const total: number = this.statusCodes.reduce((acc, v) => acc + v.hits, 0);
            return this.statusCodes.map(v => {
                return {
                    data: [
                        textService.getHttpStatusText(v.name),
                        textService.humanizeNumber(v.hits),
                    ],
                    fraction: total ? v.hits / total : 0,
                };
            });
        },
    },

    methods: {
        clickRow(rowIndex: number): void {
            this.$emit('filter', FilterDimension.Status, this.statusCodes[rowIndex].name);
        },
    },
});
