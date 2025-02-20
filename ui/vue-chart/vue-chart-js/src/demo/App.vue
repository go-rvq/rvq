<template>
  <v-app>
    <v-row>
      <v-col>
        <chart-line :data="line" :options="options" />
      </v-col>
      <v-col>
        <chart-line :data="lineRandom" :options="options" />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <chart-bar :data="bar" :options="options" />
      </v-col>
      <v-col>
        <chart-bubble :data="bubble" :options="options" />
      </v-col>
      <v-col>
        <chart-doughnut :data="doughnut" :options="options" />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <chart-pie :data="pie" :options="options" />
      </v-col>
      <v-col>
        <chart-polar-area :data="polarArea" :options="options" />
      </v-col>
      <v-col>
        <chart-radar :data="radar" :options="options" />
      </v-col>
      <v-col>
        <chart-scatter :data="scatter" :options="options" />
      </v-col>
    </v-row>
  </v-app>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import type { ChartData } from "chart.js";

const options = {
    responsive: true,
    maintainAspectRatio: false,
  },
  bar = {
    labels: ["January", "February", "March"],
    datasets: [{ data: [40, 20, 12] }],
  },
  bubble = {
    datasets: [
      {
        label: "Data One",
        backgroundColor: "#f87979",
        data: [
          {
            x: 20,
            y: 25,
            r: 5,
          },
          {
            x: 40,
            y: 10,
            r: 10,
          },
          {
            x: 30,
            y: 22,
            r: 30,
          },
        ],
      },
      {
        label: "Data Two",
        backgroundColor: "#7C8CF8",
        data: [
          {
            x: 10,
            y: 30,
            r: 15,
          },
          {
            x: 20,
            y: 20,
            r: 10,
          },
          {
            x: 15,
            y: 8,
            r: 30,
          },
        ],
      },
    ],
  },
  doughnut = {
    labels: ["VueJs", "EmberJs", "ReactJs", "AngularJs"],
    datasets: [
      {
        backgroundColor: ["#41B883", "#E46651", "#00D8FF", "#DD1B16"],
        data: [40, 20, 80, 10],
      },
    ],
  },
  line = {
    labels: ["January", "February", "March", "April", "May", "June", "July"],
    datasets: [
      {
        label: "Data One",
        backgroundColor: "#f87979",
        data: [40, 39, 10, 40, 39, 80, 40],
      },
    ],
  },
  pie = {
    labels: ["VueJs", "EmberJs", "ReactJs", "AngularJs"],
    datasets: [
      {
        backgroundColor: ["#41B883", "#E46651", "#00D8FF", "#DD1B16"],
        data: [40, 20, 80, 10],
      },
    ],
  },
  polarArea = {
    labels: [
      "Eating",
      "Drinking",
      "Sleeping",
      "Designing",
      "Coding",
      "Cycling",
      "Running",
    ],
    datasets: [
      {
        label: "My First dataset",
        backgroundColor: "rgba(179,181,198,0.2)",
        pointBackgroundColor: "rgba(179,181,198,1)",
        pointBorderColor: "#fff",
        pointHoverBackgroundColor: "#fff",
        pointHoverBorderColor: "rgba(179,181,198,1)",
        data: [65, 59, 90, 81, 56, 55, 40],
      },
      {
        label: "My Second dataset",
        backgroundColor: "rgba(255,99,132,0.2)",
        pointBackgroundColor: "rgba(255,99,132,1)",
        pointBorderColor: "#fff",
        pointHoverBackgroundColor: "#fff",
        pointHoverBorderColor: "rgba(255,99,132,1)",
        data: [28, 48, 40, 19, 96, 27, 100],
      },
    ],
  },
  radar = {
    labels: [
      "Eating",
      "Drinking",
      "Sleeping",
      "Designing",
      "Coding",
      "Cycling",
      "Running",
    ],
    datasets: [
      {
        label: "My First dataset",
        backgroundColor: "rgba(179,181,198,0.2)",
        borderColor: "rgba(179,181,198,1)",
        pointBackgroundColor: "rgba(179,181,198,1)",
        pointBorderColor: "#fff",
        pointHoverBackgroundColor: "#fff",
        pointHoverBorderColor: "rgba(179,181,198,1)",
        data: [65, 59, 90, 81, 56, 55, 40],
      },
      {
        label: "My Second dataset",
        backgroundColor: "rgba(255,99,132,0.2)",
        borderColor: "rgba(255,99,132,1)",
        pointBackgroundColor: "rgba(255,99,132,1)",
        pointBorderColor: "#fff",
        pointHoverBackgroundColor: "#fff",
        pointHoverBorderColor: "rgba(255,99,132,1)",
        data: [28, 48, 40, 19, 96, 27, 100],
      },
    ],
  },
  scatter = {
    datasets: [
      {
        label: "Scatter Dataset 1",
        fill: false,
        borderColor: "#f87979",
        backgroundColor: "#f87979",
        data: [
          {
            x: -2,
            y: 4,
          },
          {
            x: -1,
            y: 1,
          },
          {
            x: 0,
            y: 0,
          },
          {
            x: 1,
            y: 1,
          },
          {
            x: 2,
            y: 4,
          },
        ],
      },
      {
        label: "Scatter Dataset 2",
        fill: false,
        borderColor: "#7acbf9",
        backgroundColor: "#7acbf9",
        data: [
          {
            x: -2,
            y: -4,
          },
          {
            x: -1,
            y: -1,
          },
          {
            x: 0,
            y: 1,
          },
          {
            x: 1,
            y: -1,
          },
          {
            x: 2,
            y: -4,
          },
        ],
      },
    ],
  },
  lineRandom = ref<ChartData<"bar">>({
    datasets: [],
  }),
  getRandomInt = () => Math.floor(Math.random() * (50 - 5 + 1)) + 5,
  randomData = () => ({
    labels: [
      "January" + getRandomInt(),
      "February",
      "March",
      "April",
      "May",
      "June",
      "July",
      "August",
      "September",
      "October",
      "November",
      "December",
    ],
    datasets: [
      {
        label: "Data One",
        backgroundColor: "#f87979",
        data: [
          getRandomInt(),
          getRandomInt(),
          getRandomInt(),
          getRandomInt(),
          getRandomInt(),
          getRandomInt(),
          getRandomInt(),
          getRandomInt(),
          getRandomInt(),
          getRandomInt(),
          getRandomInt(),
          getRandomInt(),
        ],
      },
    ],
  });

onMounted(() => {
  setInterval(() => {
    lineRandom.value = randomData();
  }, 3000);
});
</script>
