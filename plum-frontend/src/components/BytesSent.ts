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
    name: 'BytesSent',

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
                        label: 'Resource',
                        width: null,
                        align: Align.Left,
                    },
                    {
                        label: 'Bytes sent',
                        width: '100px',
                        align: Align.Right,
                    },
                ],
            };
        },

        resources(): NamedMetrics[] {
            return metricsService.group(this.data, v => v.uris)
                .sort((a, b) => a.bytes < b.bytes ? 1 : -1);
        },

        rows(): TableRow[] {
            const total: number = this.resources.reduce((acc, v) => acc + v.bytes, 0);
            return this.resources.map(v => {
                return {
                    data: [
                        v.name,
                        textService.humanizeBytes(v.bytes),
                    ],
                    fraction: total ? v.bytes / total : 0,
                };
            });
        },
    },

    methods: {
        clickRow(rowIndex: number): void {
            this.$emit('filter', FilterDimension.Uri, this.resources[rowIndex].name);
        },
    },
});
