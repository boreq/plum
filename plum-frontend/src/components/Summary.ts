import { defineComponent, type PropType } from 'vue';
import type { RangeData } from '@/dto/Data';
import { DataService } from '@/services/DataService';
import { TextService } from '@/services/TextService';

const dataService = new DataService();
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
            const sum = this.data.reduce((acc, v) => {
                return acc + dataService.getHits(v.data);
            }, 0);
            return textService.humanizeNumber(sum);
        },

        visits(): string {
            const sum = this.data.reduce((acc, v) => {
                return acc + dataService.getVisits(v.data);
            }, 0);
            return textService.humanizeNumber(sum);
        },

        bytesSent(): string {
            const sum = this.data.reduce((acc, v) => {
                return acc + dataService.getBytesSent(v.data);
            }, 0);
            return textService.humanizeBytes(sum, 0);
        },
    },
});
