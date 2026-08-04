import { defineComponent, type PropType } from 'vue';
import type { RangeData } from '@/dto/Data';
import { MetricsService } from '@/services/MetricsService';
import { TextService } from '@/services/TextService';

const metricsService = new MetricsService();
const textService = new TextService();

export default defineComponent({
    name: 'Summary',

    props: {
        data: {
            type: Array as PropType<RangeData[]>,
            default: () => [],
        },
    },

    computed: {
        hits(): string {
            return textService.humanizeNumber(metricsService.total(this.data, v => v.hits));
        },

        visits(): string {
            return textService.humanizeNumber(metricsService.total(this.data, v => v.visits));
        },

        bytesSent(): string {
            return textService.humanizeBytes(metricsService.total(this.data, v => v.bytes), 0);
        },
    },
});
