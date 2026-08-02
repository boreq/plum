import { Component, Prop, Vue } from 'vue-property-decorator';
import { TableHeader, TableRow } from '@/dto/Table';
import Table from '@/components/Table.vue';

@Component({
    components: {
        Table,
    },
})
export default class TablePopup extends Vue {

    @Prop()
    title: string;

    @Prop()
    header: TableHeader;

    @Prop()
    rows: TableRow[];

    @Prop({default: 10})
    perPage: number;

    close(): void {
        this.$emit('close');
    }

}
