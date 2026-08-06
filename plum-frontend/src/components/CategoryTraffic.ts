import { defineComponent, type PropType } from 'vue';
import { Category, CategoryLabels } from '@/dto/Category';
import type { Data, Metrics, RangeData } from '@/dto/Data';
import { MetricsService } from '@/services/MetricsService';
import { TextService } from '@/services/TextService';

const metricsService = new MetricsService();
const textService = new TextService();

const emptyMetrics: Metrics = {
    visits: 0,
    hits: 0,
    bytes: 0,
};

export default defineComponent({
    name: 'CategoryTraffic',

    props: {
        data: {
            type: Array as PropType<RangeData[]>,
            default: () => [],
        },
        category: {
            type: String as PropType<Category>,
            required: true,
        },
        checked: {
            type: Boolean,
            default: true,
        },
        active: {
            type: Boolean,
            default: false,
        },
    },

    emits: ['filter'],

    computed: {
        title(): string {
            return CategoryLabels[this.category];
        },

        hits(): string {
            return textService.humanizeNumber(this.total(v => v.hits));
        },

        visits(): string {
            return textService.humanizeNumber(this.total(v => v.visits));
        },

        bytesSent(): string {
            return textService.humanizeBytes(this.total(v => v.bytes), 0);
        },

        hitsShare(): string {
            return this.share(v => v.hits);
        },

        visitsShare(): string {
            return this.share(v => v.visits);
        },

        bytesSentShare(): string {
            return this.share(v => v.bytes);
        },
    },

    methods: {
        total(selector: (metrics: Metrics) => number): number {
            return metricsService.total(this.data, (data: Data) => selector(this.metrics(data)));
        },

        metrics(data: Data): Metrics {
            return data.categories ? (data.categories[this.category] || emptyMetrics) : emptyMetrics;
        },

        share(selector: (metrics: Metrics) => number): string {
            const total = metricsService.total(this.data, (data: Data) => this.allCategories(data, selector));
            if (!total) {
                return '0%';
            }
            return Math.round(100 * this.total(selector) / total) + '%';
        },

        allCategories(data: Data, selector: (metrics: Metrics) => number): number {
            if (!data.categories) {
                return 0;
            }
            return Object.values(data.categories).reduce((acc, metrics) => acc + selector(metrics), 0);
        },

        click(): void {
            this.$emit('filter', this.category);
        },
    },
});
