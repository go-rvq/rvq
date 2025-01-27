<script setup lang="ts">
import { ref } from 'vue'
import { default as TreeDataTable, PostLoadEvent } from '@/lib/TreeDataTable.vue'

const value = [
  {
    ID: 'eaf4e36e-4476-4d45-a0bf-c8c386bee2df',
    Nome: 'IPC',
    PrecoAtual: 0,
    PrecoAnterior: 0,
    Variacao: 0,
    Peso: 0.95,
    IPC: 1.12889,
    Impacto: 1.07244,
    PctSobreItem: 0,
    PCtSobreSubGrupo: 0,
    PctSobreGrupo: 0,
    PctSobreIPC: 0,
    Indice: 0,
    IndiceAcumulado: 0,
    Children: [
      {
        ID: 'e00814bf-025b-43e7-8163-160041a7d06a',
        Codigo: '1.0.00.00',
        Nome: 'ALIMENTAÇÃO',
        PrecoAtual: 0,
        PrecoAnterior: 0,
        Variacao: 0,
        Peso: 0.3095,
        IPC: 3.71493,
        Impacto: 1.14977,
        PctSobreItem: 0,
        PCtSobreSubGrupo: 0,
        PctSobreGrupo: 0,
        PctSobreIPC: 0,
        Indice: 0,
        IndiceAcumulado: 0,
        Children: [
          {
            ID: '47564a63-96c4-4b49-877f-949e8c84faa1',
            Codigo: '1.1.00.00',
            Nome: 'ALIMENTAÇÃO NO DOMICÍLIO',
            PrecoAtual: 0,
            PrecoAnterior: 0,
            Variacao: 0,
            Peso: 0.2595,
            IPC: 3.36367,
            Impacto: 0.87287,
            PctSobreItem: 0,
            PCtSobreSubGrupo: 0,
            PctSobreGrupo: 0,
            PctSobreIPC: 0,
            Indice: 0,
            IndiceAcumulado: 0,
            Children: [
              {
                ID: '1ad2362d-a573-46e7-9f3f-5b89cdf06376',
                Codigo: '1.1.01.00',
                Nome: 'CEREAIS, LEGUMINOSAS E OLEAGINOSAS',
                PrecoAtual: 0,
                PrecoAnterior: 0,
                Variacao: 0,
                Peso: 0.2595,
                IPC: 3.36367,
                Impacto: 0.87287,
                PctSobreItem: 0,
                PCtSobreSubGrupo: 0,
                PctSobreGrupo: 0,
                PctSobreIPC: 0,
                Indice: 0,
                IndiceAcumulado: 0,
                Children: [
                  {
                    ID: '40901ec6-498e-4aca-a252-f5f046ea1ef6',
                    Codigo: '1.1.01.01',
                    Nome: 'Arroz tipo 1',
                    PrecoAtual: 37,
                    PrecoAnterior: 33.52667,
                    Variacao: 0.10359,
                    Peso: 0.1505,
                    IPC: 0,
                    Impacto: 1.55902,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  },
                  {
                    ID: 'f4b33f97-4013-4cac-b35b-50fb99f2febc',
                    Codigo: '1.1.01.03',
                    Nome: 'Feijão vermelho',
                    PrecoAtual: 9,
                    PrecoAnterior: 9.60667,
                    Variacao: -0.06315,
                    Peso: 0.1,
                    IPC: 0,
                    Impacto: -0.6315,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  },
                  {
                    ID: 'bd03fbf4-c038-4295-bc2f-0988e61b38f3',
                    Codigo: '1.1.01.07',
                    Nome: 'Canjica branca',
                    PrecoAtual: 5,
                    PrecoAnterior: 5.32333,
                    Variacao: -0.06073,
                    Peso: 0.009,
                    IPC: 0,
                    Impacto: -0.05465,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  }
                ]
              }
            ]
          },
          {
            ID: '91649201-3c1c-445b-805f-1821c36feeba',
            Codigo: '1.2.00.00',
            Nome: 'ALIMENTAÇÃO FORA DO DOMÍCILIO',
            PrecoAtual: 0,
            PrecoAnterior: 0,
            Variacao: 0,
            Peso: 0.05,
            IPC: 5.538,
            Impacto: 0.2769,
            PctSobreItem: 0,
            PCtSobreSubGrupo: 0,
            PctSobreGrupo: 0,
            PctSobreIPC: 0,
            Indice: 0,
            IndiceAcumulado: 0,
            Children: [
              {
                ID: 'c0edcee8-e8b6-4402-8d86-89b033fedfc2',
                Codigo: '1.2.01.00',
                Nome: 'REFEIÇÃO',
                PrecoAtual: 0,
                PrecoAnterior: 0,
                Variacao: 0,
                Peso: 0.05,
                IPC: 5.538,
                Impacto: 0.2769,
                PctSobreItem: 0,
                PCtSobreSubGrupo: 0,
                PctSobreGrupo: 0,
                PctSobreIPC: 0,
                Indice: 0,
                IndiceAcumulado: 0,
                Children: [
                  {
                    ID: '2d92515c-2dbf-4f1d-b5ed-32f4f8424357',
                    Codigo: '1.2.01.01',
                    Nome: 'Self-service',
                    PrecoAtual: 47,
                    PrecoAnterior: 44.53333,
                    Variacao: 0.05538,
                    Peso: 0.05,
                    IPC: 0,
                    Impacto: 0.2769,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        ID: '4d1599c0-0b4c-4c13-a67b-4ea3e824c85c',
        Codigo: '2.0.00.00',
        Nome: 'VESTUÁRIO',
        PrecoAtual: 0,
        PrecoAnterior: 0,
        Variacao: 0,
        Peso: 0.0155,
        IPC: 8.26,
        Impacto: 0.12803,
        PctSobreItem: 0,
        PCtSobreSubGrupo: 0,
        PctSobreGrupo: 0,
        PctSobreIPC: 0,
        Indice: 0,
        IndiceAcumulado: 0,
        Children: [
          {
            ID: 'a517da82-b4f3-4a30-abbb-dc15ac9803d4',
            Codigo: '2.1.00.00',
            Nome: 'ARTIGOS DE CONFECÇÃO',
            PrecoAtual: 0,
            PrecoAnterior: 0,
            Variacao: 0,
            Peso: 0.0155,
            IPC: 8.26,
            Impacto: 0.12803,
            PctSobreItem: 0,
            PCtSobreSubGrupo: 0,
            PctSobreGrupo: 0,
            PctSobreIPC: 0,
            Indice: 0,
            IndiceAcumulado: 0,
            Children: [
              {
                ID: '17a0fa64-4a64-4213-87b0-0058bd510c19',
                Codigo: '2.1.02.00',
                Nome: 'ROUPAS FEMININAS',
                PrecoAtual: 0,
                PrecoAnterior: 0,
                Variacao: 0,
                Peso: 0.0155,
                IPC: 8.26,
                Impacto: 0.12803,
                PctSobreItem: 0,
                PCtSobreSubGrupo: 0,
                PctSobreGrupo: 0,
                PctSobreIPC: 0,
                Indice: 0,
                IndiceAcumulado: 0,
                Children: [
                  {
                    ID: 'a0d3d9c8-04f0-4c03-bfb4-872c8b2d2e97',
                    Codigo: '2.1.02.02',
                    Nome: 'Blusa malha – manga longa – fem. adulto',
                    PrecoAtual: 55,
                    PrecoAnterior: 50.80354,
                    Variacao: 0.0826,
                    Peso: 0.0155,
                    IPC: 0,
                    Impacto: 0.12803,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        ID: 'ec30ca3d-f26c-4a31-ab02-ee1f60d51177',
        Codigo: '3.0.00.00',
        Nome: 'HABITAÇÃO',
        PrecoAtual: 0,
        PrecoAnterior: 0,
        Variacao: 0,
        Peso: 0.225,
        IPC: -0.91266,
        Impacto: -0.20534,
        PctSobreItem: 0,
        PCtSobreSubGrupo: 0,
        PctSobreGrupo: 0,
        PctSobreIPC: 0,
        Indice: 0,
        IndiceAcumulado: 0,
        Children: [
          {
            ID: '94cf4520-f6e2-4767-b543-8f6ea3b71d84',
            Codigo: '3.1.00.00',
            Nome: 'MORADIA',
            PrecoAtual: 0,
            PrecoAnterior: 0,
            Variacao: 0,
            Peso: 0.225,
            IPC: -0.91266,
            Impacto: -0.20534,
            PctSobreItem: 0,
            PCtSobreSubGrupo: 0,
            PctSobreGrupo: 0,
            PctSobreIPC: 0,
            Indice: 0,
            IndiceAcumulado: 0,
            Children: [
              {
                ID: 'da250227-d989-4d2c-b7c0-a34aeb8b41e0',
                Codigo: '3.1.01.00',
                Nome: 'DESPESAS COM MANUTENÇÃO DE CASA',
                PrecoAtual: 0,
                PrecoAnterior: 0,
                Variacao: 0,
                Peso: 0.225,
                IPC: -0.91266,
                Impacto: -0.20534,
                PctSobreItem: 0,
                PCtSobreSubGrupo: 0,
                PctSobreGrupo: 0,
                PctSobreIPC: 0,
                Indice: 0,
                IndiceAcumulado: 0,
                Children: [
                  {
                    ID: '63f44152-9039-446b-ac6e-96ba5d434d61',
                    Codigo: '3.1.01.01',
                    Nome: 'Aluguel residencial',
                    PrecoAtual: 1200,
                    PrecoAnterior: 1216.66667,
                    Variacao: -0.01369,
                    Peso: 0.15,
                    IPC: 0,
                    Impacto: -0.20535,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  },
                  {
                    ID: 'e127862f-cb66-4fc2-adfe-bdd4668e91fe',
                    Codigo: '3.1.01.03',
                    Nome: 'Energia elétrica residencial',
                    PrecoAtual: 73.67,
                    PrecoAnterior: 73.67,
                    Variacao: 0,
                    Peso: 0.075,
                    IPC: 0,
                    Impacto: 0,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        ID: '37465f95-19c0-48ef-b305-8aaf933fdcac',
        Codigo: '5.0.00.00',
        Nome: 'TRANSPORTE',
        PrecoAtual: 0,
        PrecoAnterior: 0,
        Variacao: 0,
        Peso: 0.1,
        IPC: 0,
        Impacto: 0,
        PctSobreItem: 0,
        PCtSobreSubGrupo: 0,
        PctSobreGrupo: 0,
        PctSobreIPC: 0,
        Indice: 0,
        IndiceAcumulado: 0,
        Children: [
          {
            ID: 'faa9c5f4-3a36-4699-941f-469cab62ebe0',
            Codigo: '5.2.00.00',
            Nome: 'TRANSPORTE PARTICULAR',
            PrecoAtual: 0,
            PrecoAnterior: 0,
            Variacao: 0,
            Peso: 0.1,
            IPC: 0,
            Impacto: 0,
            PctSobreItem: 0,
            PCtSobreSubGrupo: 0,
            PctSobreGrupo: 0,
            PctSobreIPC: 0,
            Indice: 0,
            IndiceAcumulado: 0,
            Children: [
              {
                ID: '7623e0dd-e0ba-4723-81a9-b2f39a43f6a6',
                Codigo: '5.2.02.00',
                Nome: 'TAXAS E COMBUSTÍVEIS',
                PrecoAtual: 0,
                PrecoAnterior: 0,
                Variacao: 0,
                Peso: 0.1,
                IPC: 0,
                Impacto: 0,
                PctSobreItem: 0,
                PCtSobreSubGrupo: 0,
                PctSobreGrupo: 0,
                PctSobreIPC: 0,
                Indice: 0,
                IndiceAcumulado: 0,
                Children: [
                  {
                    ID: '9fb4da0d-fd52-4f11-b93b-f365e73e67e5',
                    Codigo: '5.2.02.01',
                    Nome: 'Gasolina comum',
                    PrecoAtual: 6.37,
                    PrecoAnterior: 6.37,
                    Variacao: 0,
                    Peso: 0.1,
                    IPC: 0,
                    Impacto: 0,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        ID: 'b94ff706-8d01-4c92-9f83-82c308942b02',
        Codigo: '6.0.00.00',
        Nome: 'COMUNICAÇÃO',
        PrecoAtual: 0,
        PrecoAnterior: 0,
        Variacao: 0,
        Peso: 0.05,
        IPC: 0,
        Impacto: 0,
        PctSobreItem: 0,
        PCtSobreSubGrupo: 0,
        PctSobreGrupo: 0,
        PctSobreIPC: 0,
        Indice: 0,
        IndiceAcumulado: 0,
        Children: [
          {
            ID: '05fbc016-4d1d-49a9-bc42-b79ec4eacf10',
            Codigo: '6.1.00.00',
            Nome: 'TELEFONE',
            PrecoAtual: 0,
            PrecoAnterior: 0,
            Variacao: 0,
            Peso: 0.05,
            IPC: 0,
            Impacto: 0,
            PctSobreItem: 0,
            PCtSobreSubGrupo: 0,
            PctSobreGrupo: 0,
            PctSobreIPC: 0,
            Indice: 0,
            IndiceAcumulado: 0,
            Children: [
              {
                ID: 'a525955b-1649-4e61-847b-5f226f93f6f1',
                Codigo: '6.1.01.00',
                Nome: 'CELULAR',
                PrecoAtual: 0,
                PrecoAnterior: 0,
                Variacao: 0,
                Peso: 0.05,
                IPC: 0,
                Impacto: 0,
                PctSobreItem: 0,
                PCtSobreSubGrupo: 0,
                PctSobreGrupo: 0,
                PctSobreIPC: 0,
                Indice: 0,
                IndiceAcumulado: 0,
                Children: [
                  {
                    ID: 'a59f404d-ed48-483d-8631-658c481e6353',
                    Codigo: '6.1.01.01',
                    Nome: 'Conta de telefone celular',
                    PrecoAtual: 71.66667,
                    PrecoAnterior: 71.66667,
                    Variacao: 0,
                    Peso: 0.05,
                    IPC: 0,
                    Impacto: 0,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        ID: '4d2eee46-7504-48f8-b468-e5ce01bc30df',
        Codigo: '7.0.00.00',
        Nome: 'SAÚDE E CUIDADOS PESSOAIS',
        PrecoAtual: 0,
        PrecoAnterior: 0,
        Variacao: 0,
        Peso: 0.05,
        IPC: 0,
        Impacto: 0,
        PctSobreItem: 0,
        PCtSobreSubGrupo: 0,
        PctSobreGrupo: 0,
        PctSobreIPC: 0,
        Indice: 0,
        IndiceAcumulado: 0,
        Children: [
          {
            ID: '4ebedd34-f7e9-426a-a363-9e081067ca99',
            Codigo: '7.1.00.00',
            Nome: 'REMÉDIOS',
            PrecoAtual: 0,
            PrecoAnterior: 0,
            Variacao: 0,
            Peso: 0.05,
            IPC: 0,
            Impacto: 0,
            PctSobreItem: 0,
            PCtSobreSubGrupo: 0,
            PctSobreGrupo: 0,
            PctSobreIPC: 0,
            Indice: 0,
            IndiceAcumulado: 0,
            Children: [
              {
                ID: '809378d5-7e1e-41f0-84cc-954012f6e027',
                Codigo: '7.1.01.00',
                Nome: 'ANTIHIPERTENSIVO',
                PrecoAtual: 0,
                PrecoAnterior: 0,
                Variacao: 0,
                Peso: 0.05,
                IPC: 0,
                Impacto: 0,
                PctSobreItem: 0,
                PCtSobreSubGrupo: 0,
                PctSobreGrupo: 0,
                PctSobreIPC: 0,
                Indice: 0,
                IndiceAcumulado: 0,
                Children: [
                  {
                    ID: '2b3eef34-3e40-4450-bc7f-9bda7c50f4dd',
                    Codigo: '7.1.01.01',
                    Nome: 'Atenolol 50 mg; com 30 comprimidos',
                    PrecoAtual: 7.685,
                    PrecoAnterior: 7.685,
                    Variacao: 0,
                    Peso: 0.05,
                    IPC: 0,
                    Impacto: 0,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        ID: 'd811791e-00a7-43a0-89ac-38a9c802a997',
        Codigo: '8.0.00.00',
        Nome: 'EDUCAÇÃO',
        PrecoAtual: 0,
        PrecoAnterior: 0,
        Variacao: 0,
        Peso: 0.1,
        IPC: 0,
        Impacto: 0,
        PctSobreItem: 0,
        PCtSobreSubGrupo: 0,
        PctSobreGrupo: 0,
        PctSobreIPC: 0,
        Indice: 0,
        IndiceAcumulado: 0,
        Children: [
          {
            ID: 'bd1756ff-6da4-436f-8320-0a92ef30355c',
            Codigo: '8.1.00.00',
            Nome: 'GASTOS COM EDUCAÇÃO',
            PrecoAtual: 0,
            PrecoAnterior: 0,
            Variacao: 0,
            Peso: 0.1,
            IPC: 0,
            Impacto: 0,
            PctSobreItem: 0,
            PCtSobreSubGrupo: 0,
            PctSobreGrupo: 0,
            PctSobreIPC: 0,
            Indice: 0,
            IndiceAcumulado: 0,
            Children: [
              {
                ID: '7ff16ffb-2add-4d6f-acd6-32148fe37135',
                Codigo: '8.1.01.00',
                Nome: 'MENSALIDADES E TAXAS ESCOLARES',
                PrecoAtual: 0,
                PrecoAnterior: 0,
                Variacao: 0,
                Peso: 0.1,
                IPC: 0,
                Impacto: 0,
                PctSobreItem: 0,
                PCtSobreSubGrupo: 0,
                PctSobreGrupo: 0,
                PctSobreIPC: 0,
                Indice: 0,
                IndiceAcumulado: 0,
                Children: [
                  {
                    ID: '4fcb7df6-9b25-4cf9-9df0-a6f2a5d807e8',
                    Codigo: '8.1.01.04',
                    Nome: 'Ensino médio',
                    PrecoAtual: 1402.66667,
                    PrecoAnterior: 1402.66667,
                    Variacao: 0,
                    Peso: 0.1,
                    IPC: 0,
                    Impacto: 0,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        ID: '20344fc7-d8d7-41fc-ad44-447c7e805ebf',
        Codigo: '9.0.00.00',
        Nome: 'DESPESAS PESSOAIS',
        PrecoAtual: 0,
        PrecoAnterior: 0,
        Variacao: 0,
        Peso: 0.1,
        IPC: 0,
        Impacto: 0,
        PctSobreItem: 0,
        PCtSobreSubGrupo: 0,
        PctSobreGrupo: 0,
        PctSobreIPC: 0,
        Indice: 0,
        IndiceAcumulado: 0,
        Children: [
          {
            ID: '65e4935b-8430-466a-b278-a3840f086f16',
            Codigo: '9.1.00.00',
            Nome: 'SERVIÇOS PESSOAIS',
            PrecoAtual: 0,
            PrecoAnterior: 0,
            Variacao: 0,
            Peso: 0.1,
            IPC: 0,
            Impacto: 0,
            PctSobreItem: 0,
            PCtSobreSubGrupo: 0,
            PctSobreGrupo: 0,
            PctSobreIPC: 0,
            Indice: 0,
            IndiceAcumulado: 0,
            Children: [
              {
                ID: 'af818e42-ab3b-42a1-b51a-9784efc78267',
                Codigo: '9.1.01.00',
                Nome: 'GASTOS PESSOAIS DIVERSOS',
                PrecoAtual: 0,
                PrecoAnterior: 0,
                Variacao: 0,
                Peso: 0.1,
                IPC: 0,
                Impacto: 0,
                PctSobreItem: 0,
                PCtSobreSubGrupo: 0,
                PctSobreGrupo: 0,
                PctSobreIPC: 0,
                Indice: 0,
                IndiceAcumulado: 0,
                Children: [
                  {
                    ID: '824f180e-e66a-4430-9519-4fab60c30000',
                    Codigo: '9.1.01.01',
                    Nome: 'Corte de cabelo feminino',
                    PrecoAtual: 66.25,
                    PrecoAnterior: 66.25,
                    Variacao: 0,
                    Peso: 0.1,
                    IPC: 0,
                    Impacto: 0,
                    PctSobreItem: 0,
                    PCtSobreSubGrupo: 0,
                    PctSobreGrupo: 0,
                    PctSobreIPC: 0,
                    Indice: 0,
                    IndiceAcumulado: 0
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  }
]

let pl: PostLoadEvent
const count = ref(0)
</script>

<template>
  <v-btn @click="() => pl.toggle(pl.items[0])">Expand</v-btn>
  <span v-if="count">Count: {{ count }}</span>
  <TreeDataTable
    @post-load="
      (v: PostLoadEvent) => {
        pl = v
        count = pl.count()
      }
    "
    :expanded="true"
    :title="(item: any) => item.Nome"
    :children="(item: any) => item.Children"
    :headers="[{ title: 'Nome', key: 'Nome' }]"
    :items="[{ Nome: 'Root', Children: [{ Nome: 'Child', Children: [{ Nome: 'Child2' }] }] }]"
    :table-props="{ density: 'compact', itemsPerPage: count }"
  >
    <template v-slot:item.Nome="{ item, value }"> {{ value }}@ </template>
  </TreeDataTable>
</template>

<style scoped></style>
