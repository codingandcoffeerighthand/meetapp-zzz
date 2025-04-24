"use client"
import { createStore, create } from 'zustand'
import Web3 from 'web3'
import abi from '@/abi/abi.json'
import { subscribeWithSelector } from 'zustand/middleware'
import { v4 } from 'uuid'
import { useWebSocketStore } from './ws'
const useWeb3Store = create(
    subscribeWithSelector((set, get) => ({
        web3: null,
        account: "",
        chainId: null,
        contract: null,
        isConnected: false,
        isLoading: false,
        error: null,
        isAuthorized: false,
        roomId: "",
        prK: "",
        // Wallet and account
        initWeb3: async () => {
            if (typeof window !== "undefined") {
                set({ isLoading: true, error: null })
                try {
                    console.info("init web3")
                    let web3Instance;
                    const provider = new Web3.providers.WebsocketProvider(process.env.NEXT_PUBLIC_INFURA_URL);
                    web3Instance = new Web3(provider);
                    const prK = localStorage?.getItem('prK')
                    if (!prK) {
                        set({ err: new Error("Private key not found") })
                        return
                    }
                    const acc = web3Instance.eth.accounts.wallet.add(prK)
                    set({ account: acc[0].address, isConnected: true, prK, web3: web3Instance })
                } catch (error) {
                    console.error("Error initializing Web3:", error)
                    set({ error: error.message })
                } finally {
                    set({ isLoading: false })
                }
            }
        },

        addKey: async (pk) => {
            set({ isLoading: true })
            try {
                localStorage.setItem('prK', pk)
                await get().initWeb3()
            } catch (err) {
                console.error(err)
                set({ err: err })
            } finally {
                set({ isLoading: false })
            }
        },

        initContract: async () => {
            set({ isLoading: true })
            try {
                const { web3, account, checkAuth } = get()
                if (!web3 || !account) {
                    set({ error: 'Web3 not initialized or no accounts found' })
                    return
                }
                const contractInstance = new web3.eth.Contract(abi, process.env.NEXT_PUBLIC_CONTRACT_ADDRESS)
                set({ contract: contractInstance })
                await checkAuth()
            } catch (err) {
                console.error('Error initializing contract:', err)
                set({ error: err })
            } finally {
                set({ isLoading: false })
            }
        },

        getOwner: async () => {
            set({ isLoading: true })
            try {
                const { contract } = get()
                const owner = await contract.methods.owner().call()
                return owner
            }
            catch (err) {
                console.error('Error getting owner:', err)
                set({ error: err })
            }
            finally {
                set({ isLoading: false })
            }
        },
        register: async () => {
            set({ isLoading: true })
            try {
                const { account, contract } = get()
                if (!contract) {
                    set({ error: 'contract not initialized' })
                    return
                }
                await contract.methods.addAuthorized().send({ from: account })
                set({ isAuthorized: true })
            }
            catch (err) {
                console.error('Error registering:', err)
                set({ error: err })
            }
            finally {
                set({ isLoading: false })
            }
        },
        checkAuth: async () => {
            set({ isLoading: true })
            try {
                const { account, contract } = get()
                if (!contract) {
                    set({ error: 'contract not initialized' })
                    return
                }
                const isAuth = await contract.methods.checkAuthorized().call({
                    from: account
                })
                set({ isAuthorized: isAuth })
            }
            catch (err) {
                console.error('Error checking auth:', err)
                set({ error: err })
            }
            finally {
                set({ isLoading: false })
            }
        },
        createRoom: async (roomId, callback) => {
            set({ isLoading: true })
            try {
                const { account, contract } = get()
                if (!contract) {
                    set({ error: 'contract not initialized' })
                    return
                }
                if (!roomId) {
                    roomId = v4()
                }
                await contract.methods.createRoom(roomId).send({ from: account })
                set({ roomId })
                console.info(roomId)
                callback(roomId)
            }
            catch (err) {
                console.error('Error creating room:', err)
                set({ error: err })
            }
            finally {
                set({ isLoading: false })
            }
        },
    })))

useWeb3Store.subscribe(
    state => state.prK,
    async (acc) => {       // Callback receives the new value of `account`
        console.log("khoi tao contract")
        await useWeb3Store.getState().initContract()
        useWebSocketStore.getState().init(acc)
    }
)
await useWeb3Store.getState().initWeb3()

export {
    useWeb3Store
}