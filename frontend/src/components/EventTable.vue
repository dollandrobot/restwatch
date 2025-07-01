<script lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from "vue";
import type { main } from "../../wailsjs/go/models";
import { GetMessages } from "../../wailsjs/go/main/App";
import { useQuasar } from "quasar";

export default {
  setup(props, { emit }) {
    const $q = useQuasar();
    const messages = ref([]);
    const isDark = computed(() => $q.dark.isActive);

    const fetchMessages = async () => {
      try {
        const result = await GetMessages();
        messages.value = result;
      } catch (error) {
        console.log("error getting existing messages", error);
      }
    };

    const onReceiveMessage = (message: main.SimpleMessage) => {
      messages.value.unshift(message);
    };

    const onRowClick = (row) => {
      emit("row-click", row.id);
    };

    onMounted(async () => {
      await fetchMessages();
      window.runtime.EventsOn("messageReceived", onReceiveMessage);
    });

    onBeforeUnmount(() => {
      window.runtime.EventsOff("messageReceived");
    });

    return {
      messages,
      onRowClick,
      isDark,
    };
  },
};
</script>

<template>
  <q-virtual-scroll
    flat
    bordered
    type="table"
    :items="messages"
    style="max-height: 100%; height: 100%"
    :virtual-scroll-item-size="48"
    :virtual-scroll-sticky-size-start="48"
  >
    <template v-slot:before>
      <thead>
        <tr>
          <th
            :class="{ 'th-dark': isDark, 'th-light': !isDark }"
            class="text-left"
          >
            #
          </th>
          <th
            :class="{ 'th-dark': isDark, 'th-light': !isDark }"
            class="text-left"
          >
            Content
          </th>
          <th
            :class="{ 'th-dark': isDark, 'th-light': !isDark }"
            class="text-left"
          >
            Received At
          </th>
        </tr>
      </thead>
    </template>
    <template v-slot:default="{ item: row, index }">
      <tr :key="index" @click="onRowClick(row)">
        <td>#{{ index + 1 }}</td>
        <td>{{ row.content }}</td>
        <td>{{ row.receivedAt }}</td>
      </tr>
    </template>
  </q-virtual-scroll>
</template>

<style lang="sass" scoped>
.q-virtual-scroll
  overflow-y: auto

// Common header styles mixin
%th-common
  position: sticky
  top: 0
  z-index: 1
  padding: 8px 16px
  transition: all 0.3s ease
  font-weight: 500

.th-light
  @extend %th-common
  background-color: $deep-purple-1
  color: $deep-purple-10

.th-dark
  @extend %th-common
  background-color: $deep-purple-10
  color: $deep-purple-1

td
  padding: 8px 16px

tr
  cursor: pointer

  &:hover
    background-color: rgba(0, 0, 0, 0.05)
    transition: background-color 0.2s ease

// Dark mode specific hover
:deep(.body--dark) tr
  &:hover
    background-color: rgba(255, 255, 255, 0.1)
    transition: background-color 0.2s ease
</style>
