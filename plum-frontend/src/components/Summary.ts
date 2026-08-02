import { Component, Prop, Vue } from 'vue-property-decorator';
import { RangeData } from '@/dto/Data';
import { DataService } from '@/services/DataService';
import { TextService } from '@/services/TextService';

@Component
export default class Summary extends Vue {

    @Prop()
    private data: RangeData[];

    private dataService = new DataService();
    private textService = new TextService();

    get hits(): string {
        const sum = this.data.reduce((acc, v) => {
            return acc + this.dataService.getHits(v.data);
        }, 0);
        return this.textService.humanizeNumber(sum);
    }

    get visits(): string {
        const sum = this.data.reduce((acc, v) => {
            return acc + this.dataService.getVisits(v.data);
        }, 0);
        return this.textService.humanizeNumber(sum);
    }

    get bytesSent(): string {
        const sum = this.data.reduce((acc, v) => {
            return acc + this.dataService.getBytesSent(v.data);
        }, 0);
        return this.textService.humanizeBytes(sum, 0);
    }

}
