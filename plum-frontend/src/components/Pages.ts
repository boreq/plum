import { Component, Prop, Vue } from 'vue-property-decorator';
import { RangeData } from '@/dto/Data';
import { TableHeader, TableRow, Align } from '@/dto/Table';
import { UriService } from '@/services/UriService';
import { TextService } from '@/services/TextService';
import Table from '@/components/Table.vue';

class UriData {
    uri: string;
    visits: number;
    hits: number;
}

@Component({
    components: {
        Table,
    },
})
export default class Pages extends Vue {

    @Prop()
    private data: RangeData[];

    private uriService = new UriService();
    private textService = new TextService();

    get header(): TableHeader {
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
    }

    get rows(): TableRow[] {
        if (!this.data) {
            return [];
        }
        const rows: UriData[] = [];
        for (const rangeData of this.data) {
            if (rangeData.data.uris) {
                Object.entries(rangeData.data.uris).forEach(([uri, uriData]) => {
                    let row = rows.find(v => v.uri === uri);
                    if (!row) {
                        row = {
                            uri: uri,
                            visits: 0,
                            hits: 0,
                        };
                        rows.push(row);
                    }
                    row.visits += uriData.visits;
                    Object.entries(uriData.statuses).forEach(([status, statusData]) => {
                        row.hits += statusData.hits;
                    });
                });
            }
        }
        return this.toTableRows(rows);
    }

    private toTableRows(uriData: UriData[]): TableRow[] {
        const total: number = uriData.reduce((acc, v) => acc + v.visits, 0);
        return uriData
            .sort((a, b) => a.visits < b.visits ? 1 : -1)
            .filter(v => !this.uriService.isStaticResource(v.uri))
            .map(v => {
                return {
                    data: [
                        v.uri,
                        this.textService.humanizeNumber(v.hits),
                        this.textService.humanizeNumber(v.visits),
                    ],
                    fraction: v.visits / total,
                };
            });
    }

}
