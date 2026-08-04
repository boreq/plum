import { defineComponent, markRaw, type PropType } from 'vue';
import type { Dictionary, RangeData } from '@/dto/Data';
import { ChartColors } from '@/dto/ChartColors';
import { DataService } from '@/services/DataService';
import { TextService } from '@/services/TextService';
import { GroupingType } from '@/dto/GroupingType';
import { ChartAnimation } from '@/dto/ChartAnimation';
import { Chart } from '@/chartjs';

class ChartData {
    label: string;
    statuses: Dictionary<Dictionary<number>>;
}

type BarChart = Chart<'bar', number[], string>;

const dataService = new DataService();
const textService = new TextService();

export default defineComponent({
    name: 'StatusCodesChart',

    props: {
        data: {
            type: Array as PropType<RangeData[]>,
            default: () => [],
        },
        groupingType: {
            type: Number as PropType<GroupingType>,
            default: GroupingType.Hourly,
        },
    },

    emits: ['select-data'],

    data() {
        return {
            chart: null as BarChart,
        };
    },

    watch: {
        data(): void {
            this.redraw();
        },
    },

    mounted(): void {
        this.redraw();
    },

    beforeUnmount(): void {
        if (this.chart) {
            this.chart.destroy();
            this.chart = null;
        }
    },

    methods: {
        redraw(): void {
            if (!this.data) {
                return;
            }

            const chartData: ChartData[] = this.data.map(rangeData => {
                return {
                    label: textService.formatDate(rangeData.time, this.groupingType),
                    statuses: dataService.getStatusMapping(rangeData.data),
                };
            });
            this.drawChart(chartData);
        },

        drawChart(chartData: ChartData[]): void {
            const statuses = chartData.map(v => this.groupByStatusType(v));
            const labels: string[] = chartData.map(v => v.label);

            if (!this.chart) {
                this.chart = markRaw(this.createChart());
            }

            this.chart.data.labels.length = 0;
            for (const label of labels) {
                this.chart.data.labels.push(label);
            }

            // Prepare data
            const statusDatas: number[][] = [];
            for (let datasetIndex = 0; datasetIndex < this.chart.data.datasets.length; datasetIndex++) {
                statusDatas.push([]);
                statuses.forEach((statusData, statusIndex) => {
                    const total = statusData.reduce((acc, [, hits]) => {
                        return acc + hits;
                    }, 0);
                    const statusString = this.toStatusString((datasetIndex + 1).toString());
                    const element = statusData.find(v => v[0] === statusString);
                    if (element) {
                        statusDatas[datasetIndex][statusIndex] = element[1] / total;
                    } else {
                        statusDatas[datasetIndex][statusIndex] = 0;
                    }
                });
            }

            // Update with zeroes
            for (let datasetIndex = 0; datasetIndex < this.chart.data.datasets.length; datasetIndex++) {
                this.chart.data.datasets[datasetIndex].data.length = statusDatas[datasetIndex].length;
                for (let dataIndex = 0; dataIndex < statusDatas[datasetIndex].length; dataIndex++) {
                    if (!this.chart.data.datasets[datasetIndex].data[dataIndex]) {
                        this.chart.data.datasets[datasetIndex].data[dataIndex] = 0;
                    }
                }
            }
            this.chart.update('none');

            // Update with real values and animate
            for (let datasetIndex = 0; datasetIndex < this.chart.data.datasets.length; datasetIndex++) {
                statusDatas[datasetIndex].forEach((value, statusIndex) => {
                    this.chart.data.datasets[datasetIndex].data[statusIndex] = value;
                });
            }
            this.chart.update();
        },

        groupByStatusType(chartData: ChartData): Array<[string, number]> {
            return Object.entries(chartData.statuses)
                .reduce<Array<[string, number]>>((acc, [status, uriMap]) => {
                    let element = acc.find(([statusString]) => {
                        return this.toStatusString(status) === statusString;
                    });
                    if (!element) {
                        element = [this.toStatusString(status), 0];
                        acc.push(element);
                    }
                    element[1] += this.getTotal(uriMap);
                    return acc.sort((a, b) => a[0] < b[0] ? -1 : 1);
                }, []);
        },

        getTotal(uriMap: Dictionary<number>): number {
            return Object.entries(uriMap)
                .reduce((acc, [, hits]) => {
                    return acc + hits;
                }, 0);
        },

        toStatusString(status: string): string {
            return status[0] + 'xx';
        },

        createChart(): BarChart {
            return new Chart<'bar', number[], string>(this.$refs.canvas as HTMLCanvasElement, {
                type: 'bar',
                data: {
                    labels: [],
                    datasets: [
                        this.createDataset('1xx', ChartColors.Blue),
                        this.createDataset('2xx', ChartColors.Green),
                        this.createDataset('3xx', ChartColors.Violet),
                        this.createDataset('4xx', ChartColors.Orange),
                        this.createDataset('5xx', ChartColors.Red),
                    ],
                },
                options: {
                    maintainAspectRatio: false,
                    animation: {
                        duration: ChartAnimation.Duration,
                    },
                    scales: {
                        x: {
                            ticks: {
                                display: false,
                            },
                            grid: {
                                display: false,
                            },
                            stacked: true,
                        },
                        y: {
                            min: 0,
                            max: 1,
                            ticks: {
                                maxTicksLimit: 5,
                                callback: value => {
                                    return Math.round(Number(value) * 100) + '%';
                                },
                            },
                            stacked: true,
                        },
                    },
                    plugins: {
                        legend: {
                            display: false,
                        },
                        tooltip: {
                            mode: 'index',
                            callbacks: {
                                label: context => {
                                    const label = context.dataset.label;
                                    const percent = Math.round(context.parsed.y * 100);
                                    return `${label}: ${percent}%`;
                                },
                            },
                        },
                    },
                    onClick: (_event, elements) => {
                        const index = elements.length > 0 ? elements[0].index : null;
                        this.$emit('select-data', index);
                    },
                },
            });
        },

        createDataset(label: string, color: string) {
            return {
                label: label,
                data: [] as number[],
                backgroundColor: color,
                pointBackgroundColor: color,
                borderColor: color,
                borderWidth: 2,
                barPercentage: 1,
                categoryPercentage: 0.9,
            };
        },
    },
});
