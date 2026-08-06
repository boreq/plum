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
    },

    methods: {
        total(selector: (metrics: Metrics) => number): number {
            return metricsService.total(this.data, (data: Data) => selector(this.metrics(data)));
        },

        metrics(data: Data): Metrics {
            return data.categories ? (data.categories[this.category] || emptyMetrics) : emptyMetrics;
        },

        click(): void {
            this.$emit('filter', this.category);
        },
    },
});
