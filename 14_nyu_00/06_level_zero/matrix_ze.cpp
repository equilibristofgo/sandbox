#include <level_zero/ze_api.h>
#include <iostream>
#include <vector>
#include <fstream>
#include <cstring>
#include "matrix_ze.h"

#define CHECK_ZE(a) do { \
    ze_result_t res = (a); \
    if (res != ZE_RESULT_SUCCESS) { \
        std::cerr << "Level Zero Error: " << std::hex << res << " at " << #a << std::endl; \
        exit(1); \
    } \
} while (0)

static std::vector<uint8_t> read_file(const std::string& filename) {
    std::ifstream is(filename, std::ios::binary | std::ios::ate);
    std::streamsize size = is.tellg();
    is.seekg(0, std::ios::beg);
    std::vector<uint8_t> buffer(size);
    if (!is.read((char*)buffer.data(), size)) return {};
    return buffer;
}

extern "C" void matrix_multiply_ze(const float* A, const float* B, float* C, int N) {
    // 1. Initialize
    CHECK_ZE(zeInit(0));

    // 2. Get Driver
    uint32_t driverCount = 0;
    CHECK_ZE(zeDriverGet(&driverCount, nullptr));
    if (driverCount == 0) {
        std::cerr << "Error: No Level Zero drivers found." << std::endl;
        exit(1);
    }
    std::vector<ze_driver_handle_t> allDrivers(driverCount);
    CHECK_ZE(zeDriverGet(&driverCount, allDrivers.data()));
    ze_driver_handle_t driver = allDrivers[0];

    // 3. Get Device
    uint32_t deviceCount = 0;
    CHECK_ZE(zeDeviceGet(driver, &deviceCount, nullptr));
    if (deviceCount == 0) {
        std::cerr << "Error: No Level Zero devices found for driver." << std::endl;
        exit(1);
    }
    std::vector<ze_device_handle_t> allDevices(deviceCount);
    CHECK_ZE(zeDeviceGet(driver, &deviceCount, allDevices.data()));
    
    ze_device_handle_t device = nullptr;
    for(auto d : allDevices) {
        ze_device_properties_t props = {ZE_STRUCTURE_TYPE_DEVICE_PROPERTIES};
        CHECK_ZE(zeDeviceGetProperties(d, &props));
        if (props.type == ZE_DEVICE_TYPE_GPU) {
            device = d;
            std::cout << "Running on Level Zero device: " << props.name << std::endl;
            break;
        }
    }

    if (!device) {
        std::cerr << "Error: No GPU device found." << std::endl;
        exit(1);
    }

    // 4. Create Context
    ze_context_desc_t contextDesc = {ZE_STRUCTURE_TYPE_CONTEXT_DESC};
    ze_context_handle_t context;
    CHECK_ZE(zeContextCreate(driver, &contextDesc, &context));

    // 5. Create Command Queue
    uint32_t queueGroupCount = 0;
    CHECK_ZE(zeDeviceGetCommandQueueGroupProperties(device, &queueGroupCount, nullptr));
    std::vector<ze_command_queue_group_properties_t> queueGroupProps(queueGroupCount);
    for(auto& p : queueGroupProps) p.stype = ZE_STRUCTURE_TYPE_COMMAND_QUEUE_GROUP_PROPERTIES;
    CHECK_ZE(zeDeviceGetCommandQueueGroupProperties(device, &queueGroupCount, queueGroupProps.data()));

    uint32_t computeQueueGroupIndex = -1;
    for (uint32_t i = 0; i < queueGroupCount; i++) {
        if (queueGroupProps[i].flags & ZE_COMMAND_QUEUE_GROUP_PROPERTY_FLAG_COMPUTE) {
            computeQueueGroupIndex = i;
            break;
        }
    }

    ze_command_queue_desc_t queueDesc = {ZE_STRUCTURE_TYPE_COMMAND_QUEUE_DESC};
    queueDesc.ordinal = computeQueueGroupIndex;
    queueDesc.index = 0;
    queueDesc.mode = ZE_COMMAND_QUEUE_MODE_ASYNCHRONOUS;
    ze_command_queue_handle_t queue;
    CHECK_ZE(zeCommandQueueCreate(context, device, &queueDesc, &queue));

    // 6. Load Module (SPIR-V)
    auto spirv = read_file("kernel.spv_bmg.spv");
    if (spirv.empty()) {
        std::cerr << "Error: Could not read kernel.spv_bmg.spv" << std::endl;
        exit(1);
    }
    ze_module_desc_t moduleDesc = {ZE_STRUCTURE_TYPE_MODULE_DESC};
    moduleDesc.format = ZE_MODULE_FORMAT_IL_SPIRV;
    moduleDesc.inputSize = spirv.size();
    moduleDesc.pInputModule = spirv.data();
    ze_module_handle_t module;
    CHECK_ZE(zeModuleCreate(context, device, &moduleDesc, &module, nullptr));

    // 7. Create Kernel
    ze_kernel_desc_t kernelDesc = {ZE_STRUCTURE_TYPE_KERNEL_DESC};
    kernelDesc.pKernelName = "matrix_multiply";
    ze_kernel_handle_t kernel;
    CHECK_ZE(zeKernelCreate(module, &kernelDesc, &kernel));

    // 8. Allocate USM Memory
    ze_device_mem_alloc_desc_t deviceDesc = {ZE_STRUCTURE_TYPE_DEVICE_MEM_ALLOC_DESC};
    ze_host_mem_alloc_desc_t hostDesc = {ZE_STRUCTURE_TYPE_HOST_MEM_ALLOC_DESC};
    float *zeA, *zeB, *zeC;
    size_t matrixSize = N * N * sizeof(float);
    CHECK_ZE(zeMemAllocShared(context, &deviceDesc, &hostDesc, matrixSize, 1, device, (void**)&zeA));
    CHECK_ZE(zeMemAllocShared(context, &deviceDesc, &hostDesc, matrixSize, 1, device, (void**)&zeB));
    CHECK_ZE(zeMemAllocShared(context, &deviceDesc, &hostDesc, matrixSize, 1, device, (void**)&zeC));

    // Copy data to USM
    memcpy(zeA, A, matrixSize);
    memcpy(zeB, B, matrixSize);

    // 9. Set Kernel Arguments
    CHECK_ZE(zeKernelSetArgumentValue(kernel, 0, sizeof(void*), &zeA));
    CHECK_ZE(zeKernelSetArgumentValue(kernel, 1, sizeof(void*), &zeB));
    CHECK_ZE(zeKernelSetArgumentValue(kernel, 2, sizeof(void*), &zeC));
    CHECK_ZE(zeKernelSetArgumentValue(kernel, 3, sizeof(int), &N));

    // 10. Launch Kernel
    ze_command_list_desc_t listDesc = {ZE_STRUCTURE_TYPE_COMMAND_LIST_DESC};
    ze_command_list_handle_t cmdList;
    CHECK_ZE(zeCommandListCreate(context, device, &listDesc, &cmdList));

    uint32_t groupSizeX = 16, groupSizeY = 16;
    CHECK_ZE(zeKernelSetGroupSize(kernel, groupSizeX, groupSizeY, 1));
    ze_group_count_t dispatch = { (uint32_t)N / groupSizeX, (uint32_t)N / groupSizeY, 1 };

    CHECK_ZE(zeCommandListAppendLaunchKernel(cmdList, kernel, &dispatch, nullptr, 0, nullptr));
    CHECK_ZE(zeCommandListClose(cmdList));
    CHECK_ZE(zeCommandQueueExecuteCommandLists(queue, 1, &cmdList, nullptr));
    CHECK_ZE(zeCommandQueueSynchronize(queue, UINT64_MAX));

    // Copy result back
    memcpy(C, zeC, matrixSize);

    // 11. Cleanup
    zeMemFree(context, zeA);
    zeMemFree(context, zeB);
    zeMemFree(context, zeC);
    zeCommandListDestroy(cmdList);
    zeKernelDestroy(kernel);
    zeModuleDestroy(module);
    zeCommandQueueDestroy(queue);
    zeContextDestroy(context);
}
