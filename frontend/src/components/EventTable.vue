<script lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import type { main } from "../../wailsjs/go/models";
import { GetMessages } from "../../wailsjs/go/main/App";

export default {
  setup(props, { emit }) {
    const messages = ref([]);

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
          <th class="bg-primary text-left text-black no-pointer-events">#</th>
          <th class="bg-primary text-left text-black no-pointer-events">
            Content
          </th>
          <th class="bg-primary text-left text-black no-pointer-events">
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

<style scoped>
.q-virtual-scroll {
  overflow-y: auto;
}

th {
  position: sticky;
  top: 0;
  background-color: white;
  z-index: 1;
  padding: 8px 16px;
}

td {
  padding: 8px 16px;
}

tr {
  cursor: pointer;
}

tr:hover {
  background-color: rgba(0, 0, 0, 0.03);
}
</style>
