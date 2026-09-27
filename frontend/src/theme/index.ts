import type { GlobalThemeOverrides } from 'naive-ui'

function white(alpha: number): string {
  return `rgba(255, 255, 255, ${alpha})`
}

function black(alpha: number): string {
  return `rgba(0, 0, 0, ${alpha})`
}

// buildThemeOverrides creates the naive-ui theme overrides for the given mode.
export function buildThemeOverrides(dark: boolean): GlobalThemeOverrides {
  const main = dark ? '#F8F8F8' : '#171717'
  const overlay = dark ? white : black
  const border = overlay(0.1)

  return {
    common: {
      borderRadius: '8px',
      borderColor: border,
      primaryColor: main,
      primaryColorHover: dark ? '#EDEDED' : '#333333',
      primaryColorPressed: dark ? '#B8B8B8' : '#000000',
      primaryColorSuppl: main,
    },
    Layout: {
      color: dark ? '#000000' : '#FFFFFF',
      siderColor: dark ? '#000000' : '#FFFFFF',
    },
    Button: {
      textColor: overlay(0.45),
      textColorHover: main,
      textColorPressed: main,
    },
    Menu: {
      itemHeight: '36px',
      borderRadius: '8px',
      fontSize: '14px',
      fontWeightActive: '500',
      itemColorHover: overlay(0.06),
      itemColorActive: overlay(0.1),
      itemColorActiveHover: overlay(0.13),
      itemColorActiveCollapsed: overlay(0.1),
      itemTextColor: overlay(0.6),
      itemTextColorHover: overlay(0.9),
      itemTextColorActive: main,
      itemTextColorChildActive: main,
      itemIconColor: overlay(0.45),
      itemIconColorHover: overlay(0.9),
      itemIconColorActive: main,
      dividerColor: border,
    },
  }
}
