<template>
    <div>
        <Table :header="header" :rows="rows" :clickable="true" v-on:click-row="clickRow($event)"></Table>
        <TablePopup v-if="displayModal" :header="popupHeader" :rows="modalRows" :title="modalTitle"
                    v-on:close="displayModal = false"></TablePopup>
    </div>
</template>
<script lang="ts" src="./StatusCodes.ts"></script>
<style scoped lang="scss" src="./StatusCodes.scss"></style>
