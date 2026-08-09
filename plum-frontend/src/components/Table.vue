<template>
    <div class="table">
        <div class="thead" v-if="header">
            <div class="th" v-for="(column, index) in header.columns" :key="index" :style="getColumnStyle(index)"
                 v-bind:class="{sortable: isSortable(index), sorted: sortColumn === index}"
                 v-on:click="toggleSort(index)">
                {{ column.label }}<i v-if="sortColumn === index" class="sort-icon" :class="sortIconClass()"></i>
            </div>
        </div>
        <div v-if="!dataPresent" class="no-data">
            no data
        </div>
        <div class="tbody" v-if="dataPresent">
            <div class="tr" v-for="(entry, rowIndex) in limitedRows" :key="entry.index" v-on:click="click(rowIndex)" v-bind:class="{clickable: clickable}">
                <div class="td" v-for="(value, columnIndex) in entry.row.data" :key="columnIndex" :style="getColumnStyle(columnIndex)" :title="formatValue(columnIndex, value)">
                    {{ formatValue(columnIndex, value) }}<i v-if="columnIndex === 0 && entry.row.icon" class="icon" :class="entry.row.icon" :title="entry.row.iconTitle"></i>
                </div>
                <div class="background" :style="getBackgroundStyle(rowIndex)"></div>
            </div>
        </div>
        <div class="navigation" v-if="dataPresent">
            <div>
                <a v-on:click="prevPage()"
                   v-bind:class="{inactive: !hasPrevPage}">
                    <i class="fas fa-chevron-left"></i>
                </a>
            </div>
            <div class="center"></div>
            <div>
                <a v-on:click="nextPage()"
                   v-bind:class="{inactive: !hasNextPage}">
                    <i class="fas fa-chevron-right"></i>
                </a>
            </div>
        </div>
    </div>
</template>
<script lang="ts" src="./Table.ts"></script>
<style scoped lang="scss" src="./Table.scss"></style>
