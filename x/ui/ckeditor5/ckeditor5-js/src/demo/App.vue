<template>
  <v-app>
    <v-row>
      <v-col>
        <Editor v-model="data" />
      </v-col>
    </v-row>
  </v-app>
</template>

<script setup>
import { ref } from "vue";

const data = ref(`<p>hello</p>
<ol class="multilevel-list">
  <li>asdasd
    <ol>
      <li>asdasd</li>
    </ol>
  </li>
  <li>two
    <ol>
      <li>asdasd
        <ol>
          <li>
            asdasd
            <ul>
              <li>asdasd</li>
            </ul>
          </li>
        </ol>
      </li>
      <li>asdasd
        <ol>
          <li>asdasd</li>
        </ol>
      </li>
    </ol>
  </li>
</ol>`);

import Editor from "./Editor.vue";
</script>
