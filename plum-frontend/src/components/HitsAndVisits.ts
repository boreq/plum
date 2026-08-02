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
    hits: number;
    visits: number;
}

@Component
export default class HitsAndVisits extends Vue {

    @Prop()
    private data: RangeData[];

    @Prop()
    private groupingType: GroupingType;

    private chart: Chart;

    private textService = new TextService();
    private dataService = new DataService();

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
                    hits: this.dataService.getHits(rangeData.data),
                    visits: this.dataService.getVisits(rangeData.data),
                };
            });
        this.drawChart(chartData);
    }

    private drawChart(chartData: ChartData[]): void {
        const visits: number[] = chartData.map(v => v.visits);
        const hits: number[] = chartData.map(v => v.hits);
        const labels: string[] = chartData.map(v => v.label);

        if (!this.chart) {
            this.chart = this.createChart();
        }

        this.chart.data.labels.length = 0;
        for (const label of labels) {
            this.chart.data.labels.push(label);
        }

        // Update with zeroes
        this.chart.data.datasets[0].data.length = visits.length;
        visits.forEach((value, index) => {
            if (!this.chart.data.datasets[0].data[index]) {
                this.chart.data.datasets[0].data[index] = 0;
            }
        });

        this.chart.data.datasets[1].data.length = hits.length;
        hits.forEach((value, index) => {
            if (!this.chart.data.datasets[1].data[index]) {
                this.chart.data.datasets[1].data[index] = 0;
            }
        });

        this.chart.update({duration: 0});

        // Update with real values and animate
        visits.forEach((value, index) => {
            this.chart.data.datasets[0].data[index] = value;
        });

        hits.forEach((value, index) => {
            this.chart.data.datasets[1].data[index] = value;
        });

        this.chart.update({duration: ChartAnimation.Duration});
    }

    private createChart(): Chart {
        return new Chart('chart', {
            type: 'bar',
            data: {
                labels: ['a'],
                datasets: [
                    this.createDataset('Visits', ChartColors.Primary),
                    this.createDataset('Hits', ChartColors.Secondary),
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
                            callback: (value , index, values) => {
                                if (value as number >= 10000) {
                                    return value as number / 1000 + 'k';
                                }
                                return value;
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
            backgroundColor: color,
            borderColor: color,
            borderWidth: 1,
            barPercentage: 1,
            categoryPercentage: 0.9,
        };
    }
}
