import { type App } from "vue";

import {
  ArcElement,
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  Filler,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  RadialLinearScale,
  Title,
  Tooltip,
} from "chart.js";
import {
  Bar,
  Bubble,
  Doughnut,
  Line,
  Pie,
  PolarArea,
  Radar,
  Scatter,
} from "vue-chartjs";

ChartJS.register(
  Title,
  Filler,
  Tooltip,
  Legend,
  CategoryScale,
  LinearScale,
  RadialLinearScale,
  BarElement,
  PointElement,
  ArcElement,
  LineElement,
);

export const vuechart = {
  install: (app: App) => {
    app.component("chart-bar", Bar);
    app.component("chart-doughnut", Doughnut);
    app.component("chart-line", Line);
    app.component("chart-pie", Pie);
    app.component("chart-polar-area", PolarArea);
    app.component("chart-radar", Radar);
    app.component("chart-bubble", Bubble);
    app.component("chart-scatter", Scatter);
  },
};
