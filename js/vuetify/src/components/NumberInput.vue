<template>
  <v-row>
    <v-col cols="12" sm="6">
      <v-number-input
        label="Regular"
        :model-value="5"
        disabled
      ></v-number-input>
    </v-col>

    <v-col cols="12" sm="6">
      <v-number-input
        label="Regular"
        :model-value="5"
        readonly
      ></v-number-input>
    </v-col>

    <v-col cols="12" sm="6">
      <v-number-input
        label="Solo"
        :model-value="5"
        variant="solo"
        disabled
      ></v-number-input>
    </v-col>

    <v-col cols="12" sm="6">
      <v-number-input
        label="Solo"
        :model-value="5"
        variant="solo"
        readonly
      ></v-number-input>
    </v-col>

    <v-col cols="12" sm="6">
      <v-number-input
        label="Outlined"
        :model-value="5"
        variant="outlined"
        disabled
      ></v-number-input>
    </v-col>

    <v-col cols="12" sm="6">
      <v-number-input
        label="Outlined"
        :model-value="5"
        variant="outlined"
        readonly
      ></v-number-input>
    </v-col>

    <v-col cols="12" sm="6">
      <v-number-input
        label="underlined"
        :model-value="5"
        variant="underlined"
        disabled
      ></v-number-input>
    </v-col>

    <v-col cols="12" sm="6">
      <v-number-input
        label="underlined"
        :model-value="5"
        variant="underlined"
        readonly
      ></v-number-input>
    </v-col>

    <v-col cols="12" sm="6">
      <v-number-input
        :model-value="5"
        variant="outlined"
        density="compact"
        readonly
      ></v-number-input>
    </v-col>
  </v-row>
</template>
