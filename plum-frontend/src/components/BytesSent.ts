import { Component, Prop, Vue } from 'vue-property-decorator';
import { RangeData } from '@/dto/Data';
import { Align, TableHeader, TableRow } from '@/dto/Table';
import { TextService } from '@/services/TextService';
import Table from '@/components/Table.vue';

class BytesSentData {
    uri: string;
    bytes: number;
}

@Component({
    components: {
        Table,
    },
})
export default class BytesSent extends Vue {

    @Prop()
    private data: RangeData[];

    private textService = new TextService();

    get header(): TableHeader {
        return {
            columns: [
                {
                    label: 'Resource',
                    width: null,
                    align: Align.Left,
                },
                {
                    label: 'Bytes sent',
                    width: '100px',
                    align: Align.Right,
                },
            ],
        };
    }

    get rows(): TableRow[] {
        if (!this.data) {
            return [];
        }
        return this.toTableRows(this.groupBytesSentByUri(this.data));
    }

    private groupBytesSentByUri(data: RangeData[]): BytesSentData[] {
        return data
            .reduce<BytesSentData[]>((acc, rangeData) => {
                if (rangeData.data.uris) {
                    Object.entries(rangeData.data.uris)
                        .forEach(([uri, uriData]) => {
                            const row = this.findOrCreateRow(acc, uri);
                            row.bytes += uriData.bytes;
                        });
                }
                return acc;
            }, []);
    }

    private findOrCreateRow(acc: BytesSentData[], uri: string): BytesSentData {
        let row = acc.find(v => v.uri === uri);
        if (!row) {
            row = {
                uri: uri,
                bytes: 0,
            };
            acc.push(row);
        }
        return row;
    }

    private toTableRows(bytesSentData: BytesSentData[]): TableRow[] {
        const total: number = bytesSentData.reduce((acc, v) => acc + v.bytes, 0);
        return bytesSentData
            .sort((a, b) => a.bytes < b.bytes ? 1 : -1)
            .map(v => {
                return {
                    data: [
                        v.uri,
                        this.textService.humanizeBytes(v.bytes),
                    ],
                    fraction: total ? v.bytes / total : 0,
                };
            });
    }

}
