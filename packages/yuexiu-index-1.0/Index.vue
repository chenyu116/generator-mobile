<template>
	<q-layout view="lHh Lpr lFf" class="bg-white">
		<q-header
				elevated=""
				:class="
        `bg-${$store.state.global.config.theme.baseColor} ${$store.state.global.config.theme.textColor}`
      "
		>
			<q-toolbar>
				<q-toolbar-title class="text-subtitle1">
					<q-icon name="location_on" size="sm" left="" />{{
					$t($store.state.global.startPointInfo.name) || ''
					}}
				</q-toolbar-title>
				__feature.language__
			</q-toolbar>
		</q-header>
		<q-page-container class="bg">
			<IndexFeatureButton :features="features" />
			<q-page class="bg-primary"></q-page>
		</q-page-container>
	</q-layout>
</template>

<script>
	// {imports}
	import IndexFeatureButton from 'components/IndexFeatureButton.vue';
	const components = { IndexFeatureButton };
	// {endImports}
	export default {
		name: 'Index',

		components: components,

		data() {
			const self = this;
			return {
				features: [
					{
						name: '1',
						icon: '',
						onClick: function() {
							self.$router.replace('/1');
						},
					},
					{
						name: '2',
						icon: '',
					},
					{
						name: '3',
						icon: '',
					},
					{
						name: '4',
						icon: '',
					},
					{
						name: '5',
						icon: '',
					},
					{
						name: '6',
						icon: '',
					},
				],
				hasI18n: self.$i18n !== undefined,
			};
		},
		mounted() {
			console.dir(this);
		},
	};
</script>
<style scoped>
	.bg {
		width: 100%;
		height: 35em;
		top: 0;
		background: url('__data.homeBg__') no-repeat;
		background-size: cover;
	}
</style>
