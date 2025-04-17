import { create } from 'zustand'
import Web3 from 'web3'
const useWeb3Store = create((set, get) => ({
    web3: null,
    account: null,
    chainId: null,
    isConnected: false,
    isLoading: false,
    error: null,

    initWeb3: async () => {
        set({ isLoading: true, error: null })
        try {
            let web3Instance;
            if (window.ethereum) {
                web3Instance = new Web3(window.ethereum);
                window.ethereum.on("accountsChanged", (accounts) => {
                    if (accounts.length === 0) {
                        set({ account: null, isConnected: false })
                    } else {
                        get().updateAccount(accounts[0])
                    }
                })
                window.ethereum.on('chainChanged', (chainId) => {
                    set({ chainId: parseInt(chainId, 16) })
                    window.location.reload();
                });

                set({ web3: web3Instance })
                const chainId = await web3Instance.eth.getChainId();
                set({ chainId })
            } else {
                const provider = new Web3.providers.HttpProvider(process.env.NEXT_PUBLIC_INFURA_URL);
                web3Instance = new Web3(provider);
                set({ web3: web3Instance })
            }
        } catch (error) {
            console.error("Error initializing Web3:", error)
            set({ error: error.message })
        } finally {
            set({ isLoading: false })
        }
    },
    connectWallet: async () => {
        if (!get().web3) await get().initWeb3()
        set({ isLoading: true, error: null })
        try {
            const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' })
            if (!accounts || accounts.length === 0) {
                throw new Error("No accounts found")
            }
            const web3 = get().web3
            const address = accounts[0]
            set({ isConnected: true, account: address })
        } catch (error) {
            console.error("Error connecting wallet:", error)
            set({ error: error.message })
        } finally {
            set({ isLoading: false })
        }
    },

    disconnectWallet: () => {
        set({
            account: null,
            balance: null,
            isConnected: false,
        });
    }
}))

export default useWeb3Store 