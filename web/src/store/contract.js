import { create } from 'zustand'
import abi from '@/abi.json'
import useWeb3Store from './web3'

const useContractStore = create((set, get) => ({
    contract: null,
    isLoading: false,
    error: null,
    initContract: async () => {
        console.log('Initializing contract...')
        const { web3, account } = useWeb3Store.getState()
        console.log(web3, account)
        if (!web3 || !account) {
            set({ error: 'Web3 not initialized or no accounts found' })
            return
        }
        set({ isLoading: true })
        try {
            const contractInstance = new web3.eth.Contract(abi, process.env.NEXT_PUBLIC_CONTRACT_ADDRESS)
            console.log('Contract instance:', contractInstance)
            set({ contract: contractInstance })
        } catch (err) {
            console.error('Error initializing contract:', err)
            set({ error: err })
        } finally {
            set({ isLoading: false })
        }
    },
    getOwner: async () => {
        const { contract } = get()
        const owner = await contract.methods.owner().call()
        return owner
    },
    authorize: async (address) => {
        const { contract } = get()
        try {
            const rs = await contract.methods.addAuthorizedBackend(address).send({ from: address, gas: 3000000 })
            console.info(rs, address)
        } catch (err) {
            console.error('Error authorizing backend:', err)
            set({ error: err })
        }
    },
    getAuthorizedBackend: async (idx) => {
        const { contract } = get()
        const rs = await contract.methods.authorizedBackends(idx).call()
        return rs
    },
}))

export {
    useContractStore
}