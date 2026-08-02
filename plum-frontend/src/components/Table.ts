import { Component, Prop, Vue } from 'vue-property-decorator';
import { TableHeader, TableRow } from '@/dto/Table';

@Component
export default class Table extends Vue {

    @Prop()
    header: TableHeader;

    @Prop()
    rows: TableRow[];

    @Prop({default: 10})
    perPage: number;

    @Prop()
    clickable: boolean;

    page = 0;

    get dataPresent(): boolean {
        return this.rows && this.rows.length > 0;
    }

    get limitedRows(): TableRow[] {
        if (this.page > this.allPages) {
            this.page = 0;
        }
        const start = this.page * this.perPage;
        return this.rows.slice(start, start + this.perPage);
    }

    getColumnStyle(columnIndex: number): string {
        const column = this.header.columns[columnIndex];
        const styles: string[] = [];

        if (column.width) {
            styles.push(`width: ${column.width}`);
        } else {
            styles.push('flex: 1');
        }

        if (column.align) {
            styles.push(`text-align: ${column.align}`);
        }

        return styles.join(';');
    }

    getBackgroundStyle(rowIndex: number): string {
        const row = this.limitedRows[rowIndex];
        if (row.fraction) {
            const percentage = Math.round(row.fraction * 100);
            return `width: ${percentage}%;`;
        }
        return 'display: none;';
    }


    prevPage(): void {
        if (this.hasPrevPage) {
            this.page -= 1;
        }
    }

    nextPage(): void {
        if (this.hasNextPage) {
            this.page += 1;
        }
    }

    click(rowIndex: number): void {
        const index = this.perPage * this.page + rowIndex;
        this.$emit('click-row', index);
    }

    get hasNextPage() {
        return this.page < this.allPages;
    }

    get hasPrevPage() {
        return this.page > 0;
    }

    get allPages(): number {
        if (this.rows) {
            return Math.floor(this.rows.length / this.perPage);
        }
        return 0;
    }

}
