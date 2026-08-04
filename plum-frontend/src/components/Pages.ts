import { defineComponent, type PropType } from 'vue';
import type { RangeData } from '@/dto/Data';
import { Align, type TableHeader, type TableRow } from '@/dto/Table';
import { FilterDimension } from '@/dto/Filter';
import { MetricsService, type NamedMetrics } from '@/services/MetricsService';
import { UriService } from '@/services/UriService';
import { TextService } from '@/services/TextService';
import Table from '@/components/Table.vue';

const metricsService = new MetricsService();
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

    emits: ['filter'],

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

        allPages(): NamedMetrics[] {
            return metricsService.group(this.data, v => v.uris);
        },

        pages(): NamedMetrics[] {
            return Array.from(this.allPages)
                .sort((a, b) => a.visits < b.visits ? 1 : -1)
                .filter(v => !uriService.isStaticResource(v.name));
        },

        rows(): TableRow[] {
            const total: number = this.allPages.reduce((acc, v) => acc + v.visits, 0);
            return this.pages.map(v => {
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
            this.$emit('filter', FilterDimension.Uri, this.pages[rowIndex].name);
        },
    },
});
