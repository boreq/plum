import { defineComponent, type PropType } from 'vue';
import { Category, CategoryLabels } from '@/dto/Category';
import { emptyMetrics, type Data, type Metrics } from '@/dto/Data';
import { TextService } from '@/services/TextService';

const textService = new TextService();

export default defineComponent({
    name: 'CategoryTraffic',

    props: {
        data: {
            type: Object as PropType<Data>,
            default: null,
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

        metrics(): Metrics {
            return this.data?.categories?.[this.category] || emptyMetrics;
        },
    },

    methods: {
        total(selector: (metrics: Metrics) => number): number {
            return selector(this.metrics);
        },

        share(selector: (metrics: Metrics) => number): string {
            const total = this.allCategories(selector);
            if (!total) {
                return '0%';
            }
            return Math.round(100 * this.total(selector) / total) + '%';
        },

        allCategories(selector: (metrics: Metrics) => number): number {
            if (!this.data?.categories) {
                return 0;
            }
            return Object.values(this.data.categories).reduce((acc, metrics) => acc + selector(metrics), 0);
        },

        click(): void {
            this.$emit('filter', this.category);
        },
    },
});
