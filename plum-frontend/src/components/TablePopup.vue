<template>
    <div>
        <div class="modal-open modal" v-on:click="close()">
            <div class="modal-dialog" role="document">
                <div class="modal-content" v-on:click="$event.stopPropagation()">
                    <div class="modal-header">
                        {{ title }}
                    </div>
                    <div class="modal-body">
                        <Table :header="header" :rows="rows"></Table>
                    </div>
                    <div class="modal-footer">
                        <a v-on:click="close()">Close</a>
                    </div>
                </div>
            </div>
        </div>
        <div class="modal-backdrop show"></div>
    </div>
</template>
<script lang="ts" src="./TablePopup.ts"></script>
<style scoped lang="scss" src="./TablePopup.scss"></style>
