import { defineComponent, type PropType } from 'vue';
import type { RangeData } from '@/dto/Data';
import { Align, type TableHeader, type TableRow } from '@/dto/Table';
import { FilterDimension } from '@/dto/Filter';
import { browserIcon } from '@/dto/Browser';
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

        browsers(): Map<string, string> {
            const rv = new Map<string, string>();
            for (const rangeData of this.data) {
                if (!rangeData.data || !rangeData.data.userAgents) {
                    continue;
                }
                Object.entries(rangeData.data.userAgents).forEach(([name, metrics]) => {
                    if (metrics.browser) {
                        rv.set(name, metrics.browser);
                    }
                });
            }
            return rv;
        },

        rows(): TableRow[] {
            const total: number = this.userAgents.reduce((acc, v) => acc + v.visits, 0);
            return this.userAgents.map(v => {
                const browser = this.browsers.get(v.name);
                return {
                    data: [
                        v.name,
                        textService.humanizeNumber(v.hits),
                        textService.humanizeNumber(v.visits),
                    ],
                    fraction: total ? v.visits / total : 0,
                    icon: browser ? browserIcon(browser) : null,
                    iconTitle: browser ? 'Recognized as a web browser.' : null,
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
