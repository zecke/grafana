import { DataSourceSettings } from '../../types';

export interface HttpSettingsBaseProps {
  dataSourceConfig: DataSourceSettings<any, any>;
  onChange: (config: DataSourceSettings) => void;
}

export interface HttpSettingsProps extends HttpSettingsBaseProps {
  defaultUrl: string;
  showAccessOptions?: boolean;
}
