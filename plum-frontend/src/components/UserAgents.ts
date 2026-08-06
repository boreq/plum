import { defineComponent, type PropType } from 'vue';
import type { RangeData } from '@/dto/Data';
import { Align, type TableHeader, type TableRow } from '@/dto/Table';
import { FilterDimension } from '@/dto/Filter';
import { MetricsService, type NamedMetrics } from '@/services/MetricsService';
import { TextService } from '@/services/TextService';
import Table from '@/components/Table.vue';

const metricsService = new MetricsService();
const textService = new TextService();

export default defineComponent({
    name: 'UserAgents',

    components: {
        Table,
    },

    props: {
        data: {
            type: Array as PropType<RangeData[]>,
            default: () => [],
        },
    },

    emits: ['filter'],

    computed: {
        header(): TableHeader {
            return {
                columns: [
                    {
                        label: 'User agent',
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

        userAgents(): NamedMetrics[] {
            return metricsService.group(this.data, v => v.userAgents)
                .sort((a, b) => a.visits < b.visits ? 1 : -1);
        },

        rows(): TableRow[] {
            const total: number = this.userAgents.reduce((acc, v) => acc + v.visits, 0);
            return this.userAgents.map(v => {
                return {
                    data: [
                        v.name,
                        textService.humanizeNumber(v.hits),
                        textService.humanizeNumber(v.visits),
                    ],
                    fraction: total ? v.visits / total : 0,
                };
            });
        },
    },

    methods: {
        clickRow(rowIndex: number): void {
            this.$emit('filter', FilterDimension.UserAgent, this.userAgents[rowIndex].name);
        },
    },
});
