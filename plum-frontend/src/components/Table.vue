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
            <div class="tr" v-for="entry in limitedRows" :key="entry.index" v-on:click="click(entry.index)" v-bind:class="{clickable: clickable}">
                <div class="td" v-for="(value, columnIndex) in entry.row.data" :key="columnIndex" :style="getColumnStyle(columnIndex)" :title="formatValue(columnIndex, value)">
                    {{ formatValue(columnIndex, value) }}<i v-if="columnIndex === 0 && entry.row.icon" class="icon" :class="entry.row.icon" :title="entry.row.iconTitle"></i>
                </div>
                <div class="background" :style="getBackgroundStyle(entry.row)"></div>
            </div>
        </div>
        <div class="navigation" v-if="dataPresent">
            <div>
                <a class="expand" v-if="expandable" title="Click to see all data." v-on:click="expand()">
                    <i class="fas fa-expand"></i>
                </a>
            </div>
            <div class="center"></div>
            <div>
                <a v-on:click="prevPage()"
                   v-bind:class="{inactive: !hasPrevPage}">
                    <i class="fas fa-chevron-left"></i>
                </a>
                <a v-on:click="nextPage()"
                   v-bind:class="{inactive: !hasNextPage}">
                    <i class="fas fa-chevron-right"></i>
                </a>
            </div>
        </div>

        <Teleport to="body">
            <Transition name="sidebar">
                <div class="sidebar-overlay" v-if="expanded" v-on:click.self="collapse()">
                    <div class="sidebar">
                        <div class="sidebar-header">
                            <div class="sidebar-search">
                                <i class="fas fa-search"></i>
                                <input ref="search" type="text" placeholder="Search" v-model="query">
                            </div>
                            <a class="sidebar-close" title="Click to close." v-on:click="collapse()">
                                <i class="fas fa-times"></i>
                            </a>
                        </div>
                        <div class="sidebar-body">
                            <div class="table">
                                <div class="thead" v-if="header">
                                    <div class="th" v-for="(column, index) in header.columns" :key="index" :style="getColumnStyle(index)"
                                         v-bind:class="{sortable: isSortable(index), sorted: sortColumn === index}"
                                         v-on:click="toggleSort(index)">
                                        {{ column.label }}<i v-if="sortColumn === index" class="sort-icon" :class="sortIconClass()"></i>
                                    </div>
                                </div>
                                <div v-if="!matchingRows.length" class="no-data">
                                    no data
                                </div>
                                <div class="tbody">
                                    <div class="tr" v-for="entry in matchingRows" :key="entry.index" v-on:click="clickExpanded(entry.index)" v-bind:class="{clickable: clickable}">
                                        <div class="td" v-for="(value, columnIndex) in entry.row.data" :key="columnIndex" :style="getColumnStyle(columnIndex)" :title="formatValue(columnIndex, value)">
                                            {{ formatValue(columnIndex, value) }}<i v-if="columnIndex === 0 && entry.row.icon" class="icon" :class="entry.row.icon" :title="entry.row.iconTitle"></i>
                                        </div>
                                        <div class="background" :style="getBackgroundStyle(entry.row)"></div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </Transition>
        </Teleport>
    </div>
</template>
<script lang="ts" src="./Table.ts"></script>
<style scoped lang="scss" src="./Table.scss"></style>
