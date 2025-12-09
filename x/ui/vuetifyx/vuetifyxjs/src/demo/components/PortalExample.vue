<script setup lang="ts">

import { type ContentType, PortalContent } from '@/lib/Portal/types'
import Portal from '@/lib/Portal/Portal.vue'
import { ref, shallowRef } from 'vue'
import { componentByTemplate } from '@/lib/component-by-template'

import 'vue-prism-editor/dist/prismeditor.min.css'
import { languages } from 'prismjs/components'
import { highlight } from 'prismjs'
import { PrismEditor } from 'vue-prism-editor'
import 'prismjs/components/prism-markup'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-clike'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-typescript'
import { VCard, VBtn } from 'vuetify/components'
import 'prismjs/themes/prism-coy.css'

const templateContent = ref(`<h1>My dynamic portal content</h1>
<p><b>X = </b> <code>{{x}}</code></p>
<p><b>The Message is:</b> <code>{{message}}</code></p>
100
<component :is="MyComponent" />`),
  dataContent = ref(`{
  "x": 100
}`)

const myComponentContentValue = '<p>My Component!</p>'
const myComponentContent = ref(myComponentContentValue)

const MyComponent = shallowRef(componentByTemplate({
    template: myComponentContentValue
  })),
  message = ref('single message value'),
  data = {
    message,
    MyComponent
  }

const content = ref<ContentType>(``)
const loading = ref(false)

const update = () => {
  MyComponent.value = componentByTemplate({
    template: myComponentContent.value
  })
  content.value = {
    template: templateContent.value,
    options: {
      data: (dataContent.value ? JSON.parse(dataContent.value) : {}) as Record<string, any>
    }
  } as PortalContent
}

const empty = () => {
  content.value = ''
}

const portalDeclarations = [`<Portal ... :data="data">
    <template v-pre>
        <h1>Default portal slot value</h1>
        <p><b>Message:</b> <code>{{message}}</code></p>
        <hr />
        <component :is="MyComponent" />

        <!-- bad! 'user' var is not accessible! only vars in 'data' is accessible. -->
        <p>User name: {{user}}</p>
    </template>
</Portal>`,
  `<Portal ... :data="data">
    <template #default="{message, MyComponent}">
        <h1>Default portal slot value</h1>
        <p><b>Message:</b> <code>{{message}}</code></p>
        <hr />
        <component :is="MyComponent.value" />

        <!-- Good! all vars in 'data' (by destructuring slot props) and other is accessible. -->
        <p>User name: {{user}}</p>
    </template>
</Portal>`
]
</script>
<template>
  <h1>Portal</h1>
  <VRow>
    <VCol>
      <VCard title="Usage">
        <VCardText style="overflow-y: auto; max-height: 440px">
        <p>
          Considere, in this example:
        </p>
          <PrismEditor readonly line-numbers :model-value="`const myComponentContent = ref(''<p>My Component!</p>''),
  MyComponent = shallowRef(componentByTemplate({
    template: myComponentContent.value
  })),
  message = ref('single message value'),
  data = {
    message,
    MyComponent
  },
  user = ref('John')`"
                       :highlight="(code:string) => highlight(code, languages.ts, 'ts')" />
          <p class="my-3">and Portal declaration:</p>
          <PrismEditor readonly line-numbers :model-value="portalDeclarations[0]"
                       :highlight="(code:string) => highlight(code, languages.markup, 'markup')" />
          <p class="my-3">
            If 'default' slot for Portal is wraped by tag <b><code>`&lt;template v-pre&gt;`</code></b>, then access
            only the portal data variables, similar to initialize <b><code>modelValue</code></b> property with string
            contents.
          </p>
          <p>
            Other wise, destructures the variables from default slot props.
          </p>
          <PrismEditor readonly line-numbers :model-value="portalDeclarations[1]"
                       :highlight="(code:string) => highlight(code, languages.markup, 'markup')" />
        </VCardText>
      </VCard>
    </VCol>
    <VCol>
      <VCard title="'message' data">
        <VCardText style="overflow-y: auto; max-height: 440px">
          <VTextarea v-model="message" />
        </VCardText>
      </VCard>
    </VCol>
  </VRow>
  <VRow>
    <VCol>
      <VCard title="Dynamic Portal Content">
        <VCardText style="overflow-y: auto; max-height: 440px">
          <PrismEditor placeholder="Put template here" line-numbers v-model="templateContent"
                       :highlight="(code:string) => highlight(code, languages.markup, 'markup')" />
        </VCardText>
      </VCard>
    </VCol>
    <VCol>
      <VCard title="'MyComponent' var Template">
        <VCardText style="overflow-y: auto; max-height: 440px">
          <PrismEditor placeholder="Put template here" line-numbers v-model="myComponentContent"
                       :highlight="(code:string) => highlight(code, languages.markup, 'markup')" />
        </VCardText>
      </VCard>
    </VCol>
    <VCol>
      <VCard title="Data">
        <VCardText style="overflow-y: auto; max-height: 440px">
          <PrismEditor line-numbers v-model="dataContent"
                       :highlight="(code:string) => highlight(code, languages.json, 'json')" />
        </VCardText>
      </VCard>
    </VCol>
  </VRow>
  <VRow>
    <VCol>
      <VCard>
        <VCardTitle>
          Portal
          <VBtn color="secondary" @click="content = null" class="ms-3" prepend-icon="mdi-restore" size="small" variant="flat">Restore to initial value</VBtn>
          <VBtn color="teal-lighten-4" @click="empty" class="ms-3" prepend-icon="mdi-layers-off" size="small" variant="flat">Empty</VBtn>
          <VBtn color="primary" @click="update" class="ms-3" prepend-icon="mdi-update" size="small" variant="flat">Set dynamic content</VBtn>
        </VCardTitle>
        <VDivider />
        <VCardText style="overflow-y: auto; max-height: 500px">
          <Portal tag="div" name="examplePortal" class="xy" v-model="content" v-model:loading="loading" :data="data">
            <template v-pre>
              <h1>Default portal slot value</h1>
              <p><b>Message:</b> <code>{{message}}</code></p>
              <hr />
              <component :is="MyComponent" />
            </template>
          </Portal>
        </VCardText>
      </VCard>
    </VCol>
  </VRow>
</template>
