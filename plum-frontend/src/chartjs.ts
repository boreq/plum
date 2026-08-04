import { Chart, registerables } from 'chart.js';

// Chart.js 3+ ships nothing registered by default, so every controller, scale
// and plugin the charts rely on has to be pulled in once.
Chart.register(...registerables);

export { Chart };
