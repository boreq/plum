import { Component, Prop, Vue, Watch } from 'vue-property-decorator';
import { RangeData } from '@/dto/Data';
import { ChartColors } from '@/dto/ChartColors';
import { DataService } from '@/services/DataService';
import { TextService } from '@/services/TextService';
import { GroupingType } from '@/dto/GroupingType';
import { ChartAnimation } from '@/dto/ChartAnimation';
import Chart from 'chart.js';

class ChartData {
    label: string;
    bytes: number;
}

@Component
export default class HitsAndVisits extends Vue {

    @Prop()
    private data: RangeData[];

    @Prop()
    private groupingType: GroupingType;

    private chart: Chart;

    private dataService = new DataService();
    private textService = new TextService();

    @Watch('data')
    onRangeDataChanged(value: RangeData[], oldValue: RangeData[]) {
        this.redraw();
    }

    mounted(): void {
        this.redraw();
    }

    private redraw(): void {
        if (!this.data) {
            return;
        }

        const chartData: ChartData[] = this.data
            .map(rangeData => {
                return {
                    label: this.textService.formatDate(rangeData.time, this.groupingType),
                    bytes: this.dataService.getBytesSent(rangeData.data),
                };
            });
        this.drawChart(chartData);
    }

    private drawChart(chartData: ChartData[]): void {
        const bytes: number[] = chartData.map(v => v.bytes);
        const labels: string[] = chartData.map(v => v.label);

        if (!this.chart) {
            this.chart = this.createChart();
        }

        this.chart.data.labels.length = 0;
        for (const label of labels) {
            this.chart.data.labels.push(label);
        }

        // Update with zeroes
        this.chart.data.datasets[0].data.length = bytes.length;
        bytes.forEach((value, i) => {
            if (!this.chart.data.datasets[0].data[i]) {
                this.chart.data.datasets[0].data[i] = 0;
            }
        });
        this.chart.update({duration: 0});

        // Update with real values and animate
        bytes.forEach((value, i) => {
            this.chart.data.datasets[0].data[i] = value;
        });
        this.chart.update({duration: ChartAnimation.Duration});
    }

    private createChart(): Chart {
        return new Chart('bytes-sent-chart', {
            type: 'line',
            data: {
                labels: ['a'],
                datasets: [
                    this.createDataset('Bytes sent', ChartColors.Primary),
                ],
            },
            options: {
                maintainAspectRatio: false,
                scales: {
                    xAxes: [{
                        ticks: {
                            display: false,
                        },
                        gridLines: {
                            display: false,
                        },
                        stacked: true,
                    }],
                    yAxes: [{
                        ticks: {
                            beginAtZero: true,
                            maxTicksLimit: 5,
                            callback: (value, index, values) => {
                                return this.textService.humanizeBytes(value as number);
                            },
                        },
                        stacked: false,
                    }],
                },
                legend: {
                    display: false,
                },
                tooltips: {
                    mode: 'index',
                    callbacks: {
                        label: (tooltipItem, data) => {
                            return this.textService.humanizeBytes(tooltipItem.yLabel as number);
                        },
                    },
                },
                onClick: evt => {
                    const element = this.chart.getElementAtEvent(evt) as any[];
                    const index = element && element.length > 0 ? element[0]._index : null;
                    this.$emit('select-data', index);
                },
            },
        });
    }

    private createDataset(label: string, color: string) {
        return {
            label: label,
            data: [],
            backgroundColor: 'rgba(0, 0, 0, 0)',
            pointBackgroundColor: color,
            borderColor: color,
            borderWidth: 2,
            barPercentage: 1,
            categoryPercentage: 0.9,
        };
    }

}
