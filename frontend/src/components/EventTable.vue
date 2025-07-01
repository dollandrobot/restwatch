<script lang="ts">
import { ref } from "vue";
import type { main } from "../../wailsjs/go/models";

const pagination = {
  sortBy: "receivedAt",
  descending: true,
};

const columns = [
  {
    name: "content",
    required: true,
    label: "Event",
    align: "left",
    field: (row) => row.content,
    format: (val) => `${val}`,
    sortable: true,
  },
  {
    name: "receivedAt",
    required: true,
    label: "Received At",
    align: "left",
    field: (row) => row.receivedAt,
    format: (val) => `${val}`,
    sortable: true,
    sortOrder: "da",
  },
];

const messages = ref([]);

const onReceiveMessage = (message: main.SimpleMessage) => {
  console.log("received", message);
  messages.value.unshift(message);
};

window.runtime.EventsOn("messageReceived", onReceiveMessage);

export default {
  setup() {
    return {
      columns,
      messages,
      pagination,
      onRowClick: (row) => alert(`${row.id} clicked`),
    };
  },
};
</script>

<template>
  <q-table
    flat
    bordered
    title="Received Events"
    :rows="messages"
    :columns="columns"
    v-model:pagination="pagination"
    no-data-label="Waiting for events..."
    row-key="id"
    class="col full-width event-table"
    :rows-per-page-options="[0]"
    virtual-scroll
    hide-pagination
  >
    <template v-slot:body="props">
      <q-tr
        class="cursor-pointer"
        :props="props"
        @click="onRowClick(props.row)"
      >
        <q-td key="content" :props="props">
          {{ props.row.content }}
        </q-td>
        <q-td key="receivedAt" :props="props">
          {{ props.row.receivedAt }}
        </q-td>
      </q-tr>
    </template>
  </q-table>
</template>

<style scoped>
.event-table {
  height: 99%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.event-table .q-table__middle {
  flex: 1;
  overflow: auto;
}

.q-table__bottom {
  position: sticky;
  bottom: 0;
  background-color: inherit;
}
</style>
